package generator_v2

import (
    "fmt"
    "github.com/boltdb/bolt"
    "github.com/dgraph-io/dgo"
    "github.com/ssor/sql2graphql/helper"
    "io/ioutil"
)

func MutationObjs(files ...string) error {
    for _, file := range files {
        logger.Pass("try to import file ", file)
        tableName, _, rows, err := ParseInsertSql(ioutil.ReadFile(file))
        if err != nil {
            return fmt.Errorf("parse insert sql failed: %s", err)
        }
        switch tableName {
        case "customers":
            customers, err := generateCustomers(tableName, rows)
            if err != nil {
                logger.Failedf("import %f failed for %s", file, err)
                return err
            }
            for _, customer := range customers {
                customer.UpdateUid(getUidByIndex)
                customer.updateDependentObjectUid(getUidByIndex)
                err := mutateObjects(customer, dgClient)
                if err != nil {
                    logger.Failedf("mutate customer [%d] failed for: %s", customer.CustomerNumber, err)
                    return err
                }
            }
        case "employees":
            employees, err := generateEmployees(tableName, rows)
            if err != nil {
                logger.Failedf("import %f failed for %s", file, err)
                return err
            }
            for _, employee := range employees {
                employee.UpdateUid(getUidByIndex)
                employee.updateDependentObjectUid(getUidByIndex)
                err := mutateObjects(employee, dgClient)
                if err != nil {
                    logger.Failedf("mutate employee [%d] failed for: %s", employee.EmployeeNumber, err)
                    continue
                }
            }
        default:
            logger.Failedf("do not support file %s to import", file)
        }
    }
    return nil
}

type UidSavable interface {
    GetUidInfo() (index string, uid string)
    SetUid(uid string)
}

func mutateObjects(object UidSavable, dgClient *dgo.Dgraph) error {
    mutatableObj, ok := object.(helper.Mutatable)
    if ok == false {
        return fmt.Errorf("object is not Mutatable")
    }
    //{
    //    index, uid := object.GetUidInfo()
    //    logger.Infof("before mutation index = [%s]  uid = [%s]", index, uid)
    //}
    uid, err := helper.MutationObj(mutatableObj, dgClient)
    if err != nil {
        logger.Failedf("mutate object failed for %s", err)
        return err
    }
    object.SetUid(uid)
    index, uid := object.GetUidInfo()
    //logger.Infof("after mutation index = [%s]  uid = [%s]", index, uid)
    err = saveIndexUidMap(index, uid)
    if err != nil {
        logger.Failed("save index failed for ", err)
        return err
    }
    //logger.Successf("mutate object %s and save to db success", index)
    return nil
}

func saveIndexUidMap(index, uid string) error {
    e := kvdb.Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists(normalBucketName)
        if err != nil {
            return fmt.Errorf("create bucket: %s", err)
        }
        err = bucket.Put([]byte(index), []byte(uid))
        if err != nil {
            logger.Failedf("failed to put index [%s] value [%s] for %s", index, uid, err)
            return err
        }
        return nil
    })
    if e != nil {
        logger.Failedf("save [%s] -> [%s] to db failed", index, uid)
        return e
    }
    //logger.Successf("save [%s] -> [%s] to db ", index, uid)
    return nil
}

func getUidByIndex(index string) (uid string, e error) {
    e = kvdb.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket(normalBucketName)
        if bucket == nil {
            logger.Failedf("cannot find bucket ")
            return fmt.Errorf("cannot find bucket")
        }
        v := bucket.Get([]byte(index))
        uid = string(v)
        return nil
    })
    return
}

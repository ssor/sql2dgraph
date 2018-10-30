package helper

import (
    "fmt"
    "github.com/boltdb/bolt"
    "github.com/dgraph-io/dgo"
)

func RemoveUidByIndex(index string) error {
    e := kvdb.Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists(normalBucketName)
        if err != nil {
            return fmt.Errorf("create bucket: %s", err)
        }
        err = bucket.Delete([]byte(index))
        if err != nil {
            logger.Failedf("failed to delete index [%s] for %s", index, err)
            return err
        }
        return nil
    })
    return e
}

func GetUidByIndex(qb Mutatable, client *dgo.Dgraph) (uid string, e error) {
    index, _ := qb.GetUidInfo()
    uid, e = getUidByIndexFromLocalDb(index)
    if e != nil {
        return
    }
    if len(uid) <= 0 {
        uid, e = QueryUid(qb, client)
        if e != nil {
            return
        }
        if len(uid) <= 0 {
            return
        }
        err := saveIndexUidMapToLocalDb(index, uid)
        if err != nil {
            return
        }
    }
    //logger.Infof("query by [%s] and get uid [%s]", index, uid)
    return
}

func getUidByIndexFromLocalDb(index string) (uid string, e error) {
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

func saveIndexUidMapToLocalDb(index, uid string) error {
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
    logger.Successf("save [%s] -> [%s] to db ", index, uid)
    return nil
}

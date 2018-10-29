package generator_v2

import (
    "fmt"
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

        logger.Pass("table name: ", tableName)
        switch tableName {
        case "customers":
            customers, err := generateCustomers(tableName, rows)
            if err != nil {
                logger.Failedf("import %f failed for %s", file, err)
                return err
            }
            for _, customer := range customers {
                if customer.Employee != nil {
                    uid, err := helper.GetUidByIndex(customer.Employee, dgClient)
                    if err != nil {
                        logger.Failedf("failed to get uid for employee %d", customer.Employee.EmployeeNumber)
                        continue
                    }
                    customer.Employee.SetUid(uid)
                }
                _, err := helper.MutationObj(customer, dgClient)
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
                if employee.ReportsTo != nil {
                    uid, err := helper.GetUidByIndex(employee.ReportsTo, dgClient)
                    if err != nil {
                        logger.Failedf("failed to get uid for employee %d", employee.ReportsTo.EmployeeNumber)
                        continue
                    }
                    employee.ReportsTo.SetUid(uid)
                }
                _, err := helper.MutationObj(employee, dgClient)
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

package generator_v2

import (
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "github.com/ssor/sql2graphql/helper"
    "github.com/ssor/zlog"
    "github.com/xwb1989/sqlparser"
    "strconv"
)

/*
CREATE TABLE employees (
  employeeNumber int(11) NOT NULL,
  lastName varchar(50) NOT NULL,
  firstName varchar(50) NOT NULL,
  extension varchar(10) NOT NULL,
  email varchar(100) NOT NULL,
  officeCode varchar(10) NOT NULL,
  reportsTo int(11) DEFAULT 0,
  jobTitle varchar(50) NOT NULL,
  PRIMARY KEY (employeeNumber),
  KEY reportsTo (reportsTo),
  KEY officeCode (officeCode)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/
type Employee struct {
    Uid            string    `json:"uid,omitempty"`
    EmployeeNumber int       `json:"employee_number,omitempty"`
    LastName       string    `json:"last_name,omitempty"`
    FirstName      string    `json:"first_name,omitempty"`
    Extension      string    `json:"extension,omitempty"`
    Email          string    `json:"email,omitempty"`
    Office         *Office   `json:"office_work_in,omitempty"`
    ReportsTo      *Employee `json:"report_to,omitempty"`
    JobTitle       string    `json:"job_title,omitempty"`
}

func newEmptyEmployee() *Employee {
    return &Employee{}
}

func newEmployee(number int, reportsTo *Employee) *Employee {
    return &Employee{
        EmployeeNumber: number,
        ReportsTo:      reportsTo,
    }
}

func (employee *Employee) Schemes() helper.Schemas {
    var schemes helper.Schemas
    schemes = schemes.Add(helper.NewSchemaIntIndex("employee_number")).
        Add(helper.NewSchemaString("job_title"))

    return schemes
}

func (employee *Employee) QueryBy() (string, string) {
    return "employee_number", strconv.FormatInt(int64(employee.EmployeeNumber), 10)
}

func (employee *Employee) DependentObjectHasUid() bool {
    if employee.ReportsTo != nil {
        if len(employee.ReportsTo.Uid) <= 0 {
            index, _ := employee.ReportsTo.GetUidInfo()
            logger.Failedf("employee %d  reportsto [%d] which is lack of uid, should has uid by index [%s]",
                employee.EmployeeNumber, employee.ReportsTo.EmployeeNumber, index)
            return false
        }
    }
    return true
}

func (employee *Employee) SetUid(uid string) {
    employee.Uid = uid
}

func (employee *Employee) GetUidInfo() (string, string) {
    return fmt.Sprintf("employee_%d", employee.EmployeeNumber), employee.Uid
}

func (employee *Employee) SetValue(index int, value interface{}) {
    switch index {
    case 0:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            i, err := strconv.Atoi(string(value.Val))
            if err != nil {
                zlog.Failedf("cannot parse %s to int", string(value.Val))
                return
            }
            employee.EmployeeNumber = i
        default:
            zlog.Failed("not supported value type for EmployeeNumber ")
            spew.Dump(value)
        }
    case 1:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            employee.LastName = string(value.Val)
        case *sqlparser.NullVal:
            employee.LastName = ""
        default:
            zlog.Failed("unknown value type for LastName ")
            spew.Dump(value)
        }
    case 2:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            employee.FirstName = string(value.Val)
        case *sqlparser.NullVal:
            employee.FirstName = ""
        default:
            zlog.Failed("unknown value type for FirstName ")
            spew.Dump(value)
        }
    case 3:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            employee.Extension = string(value.Val)
        case *sqlparser.NullVal:
            employee.Extension = ""
        default:
            zlog.Failed("unknown value type for Extension ")
            spew.Dump(value)
        }
    case 4:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            employee.Email = string(value.Val)
        case *sqlparser.NullVal:
            employee.Email = ""
        default:
            zlog.Failed("unknown value type for Email ")
            spew.Dump(value)
        }
    case 5:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            employee.Office = newOffice(string(value.Val))
        case *sqlparser.NullVal:
            employee.Office = nil
        default:
            zlog.Failed("unknown value type for Office ")
            spew.Dump(value)
        }
    case 6:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            i, err := strconv.Atoi(string(value.Val))
            if err != nil {
                zlog.Failed("not supported value type for ReportsTo ")
                spew.Dump(value)
            }
            employee.ReportsTo = newEmployee(i, nil)
        case *sqlparser.NullVal:
            employee.ReportsTo = nil
        default:
            zlog.Failed("unknown value type for ReportsTo ")
            spew.Dump(value)
        }
    case 7:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            employee.JobTitle = string(value.Val)
        case *sqlparser.NullVal:
            employee.JobTitle = ""
        default:
            zlog.Failed("unknown value type for JobTitle ")
            spew.Dump(value)
        }
    default:
        zlog.Failedf("index %d out of table Employee range", index)
    }
}

func generateEmployees(tableName string, rows sqlparser.Values) (employees []*Employee, e error) {
    if tableName != "employees" {
        e = fmt.Errorf("expect table name %s, and in fact is %s", "employees", tableName)
        zlog.Failed("parse employees values failed: ", e)
        return
    }

    for _, row := range rows {
        emptyEmployee := newEmptyEmployee()
        vt := sqlparser.ValTuple(row)
        for index, expr := range vt {
            emptyEmployee.SetValue(index, expr)
        }
        employees = append(employees, emptyEmployee)
    }
    return
}

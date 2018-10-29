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
CREATE TABLE customers (
  customerNumber int(11) NOT NULL,
  customerName varchar(50) NOT NULL,
  contactLastName varchar(50) NOT NULL,
  contactFirstName varchar(50) NOT NULL,
  phone varchar(50) NOT NULL,
  addressLine1 varchar(50) NOT NULL,
  addressLine2 varchar(50) DEFAULT NULL,
  city varchar(50) NOT NULL,
  state varchar(50) DEFAULT NULL,
  postalCode varchar(15) DEFAULT NULL,
  country varchar(50) NOT NULL,
  salesRepEmployeeNumber int(11) DEFAULT NULL,
  creditLimit decimal(10,2) DEFAULT NULL,
  PRIMARY KEY (customerNumber),
  KEY salesRepEmployeeNumber (salesRepEmployeeNumber)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/
type Customer struct {
    Uid              string    `json:"uid,omitempty"`
    CustomerNumber   int       `json:"customer_number,omitempty"`
    CustomerName     string    `json:"customer_name,omitempty"`
    ContactLastName  string    `json:"contact_last_name,omitempty"`
    ContactFirstName string    `json:"contact_first_name,omitempty"`
    Phone            string    `json:"phone,omitempty"`
    AddressLine1     string    `json:"address_line1,omitempty"`
    AddressLine2     string    `json:"address_line2,omitempty"`
    City             string    `json:"city,omitempty"`
    State            string    `json:"state,omitempty"`
    PostalCode       string    `json:"postal_code,omitempty"`
    Country          string    `json:"country,omitempty"`
    Employee         *Employee `json:"service_by_employee,omitempty"`
    CreditLimit      float64   `json:"credit_limit,omitempty"`
}

func newEmptyCustomer() *Customer {
    return &Customer{}
}

func newCustomer(number int) *Customer {
    return &Customer{
        CustomerNumber: number,
    }
}

func (customer *Customer) Schemes() helper.Schemas {
    var schemes helper.Schemas
    schemes = schemes.Add(helper.NewSchemaIntIndex("customer_number")).
        Add(helper.NewSchemaStringExactIndex("city"))

    return schemes
}

func (customer *Customer) SetUid(uid string) {
    customer.Uid = uid
}

func (customer *Customer) QueryBy() []interface{} {
    return []interface{}{"customer_number", customer.CustomerNumber}
}

func (customer *Customer) GetUidInfo() (string, string) {
    return fmt.Sprintf("employee_%d", customer.CustomerNumber), customer.Uid
}

func (customer *Customer) SetValue(index int, value interface{}) {
    switch index {
    case 0:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            i, err := strconv.Atoi(string(value.Val))
            if err != nil {
                zlog.Failedf("cannot parse %s to int", string(value.Val))
                return
            }
            customer.CustomerNumber = i
        default:
            zlog.Failed("not supported value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 1:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.CustomerName = string(value.Val)
        case *sqlparser.NullVal:
            customer.CustomerName = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 2:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.ContactLastName = string(value.Val)
        case *sqlparser.NullVal:
            customer.ContactLastName = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 3:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.ContactFirstName = string(value.Val)
        case *sqlparser.NullVal:
            customer.ContactFirstName = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 4:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.Phone = string(value.Val)
        case *sqlparser.NullVal:
            customer.Phone = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 5:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.AddressLine1 = string(value.Val)
        case *sqlparser.NullVal:
            customer.AddressLine1 = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 6:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.AddressLine2 = string(value.Val)
        case *sqlparser.NullVal:
            customer.AddressLine2 = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 7:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.City = string(value.Val)
        case *sqlparser.NullVal:
            customer.City = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 8:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.State = string(value.Val)
        case *sqlparser.NullVal:
            customer.State = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 9:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.PostalCode = string(value.Val)
        case *sqlparser.NullVal:
            customer.PostalCode = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 10:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            customer.Country = string(value.Val)
        case *sqlparser.NullVal:
            customer.Country = ""
        default:
            zlog.Failed("unknown value type for CustomerNumber ")
            spew.Dump(value)
        }
    case 11:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            i, err := strconv.Atoi(string(value.Val))
            if err != nil {
                zlog.Failed("not supported value type for salesRepEmployeeNumber ")
                spew.Dump(value)
            }
            customer.Employee = newEmployee(i, nil)
        case *sqlparser.NullVal:
            customer.Employee = nil
        default:
            zlog.Failed("unknown value type for salesRepEmployeeNumber ")
            spew.Dump(value)
        }
    case 12:
        switch value := value.(type) {
        case *sqlparser.SQLVal:
            f, err := strconv.ParseFloat(string(value.Val), 64)
            if err != nil {
                zlog.Failed("not supported value type for CreditLimit ")
                spew.Dump(value)
            }
            customer.CreditLimit = f
        case *sqlparser.NullVal:
            customer.CreditLimit = 0
        default:
            zlog.Failed("unknown value type for CreditLimit ")
            spew.Dump(value)
        }
    default:
        zlog.Failedf("index %d out of table customer range", index)
    }

}

//
//func (customer *Customer) updateDependentObjectUid(uidGetter func(string) (string, error)) error {
//    if customer.Employee == nil {
//        return nil
//    }
//
//    index, _ := customer.Employee.GetUidInfo()
//    uid, err := uidGetter(index)
//    if err != nil {
//        return err
//    }
//    if len(uid) <= 0 {
//        logger.Failedf("get uid of %s failed", index)
//    } else {
//        customer.Employee.SetUid(uid)
//    }
//    return nil
//}

func (customer *Customer) DependentObjectHasUid() bool {
    if customer.Employee != nil {
        if len(customer.Employee.Uid) <= 0 {
            return false
        }
    }
    return true
}

func generateCustomers(tableName string, rows sqlparser.Values) (customers []*Customer, e error) {

    if tableName != "customers" {
        e = fmt.Errorf("expect table name %s, and in fact is %s", "customers", tableName)
        zlog.Failed("parse customer values failed: ", e)
        return
    }

    for _, row := range rows {
        emptyCustomer := newEmptyCustomer()
        vt := sqlparser.ValTuple(row)
        for index, expr := range vt {
            emptyCustomer.SetValue(index, expr)
        }
        customers = append(customers, emptyCustomer)
    }
    return
}

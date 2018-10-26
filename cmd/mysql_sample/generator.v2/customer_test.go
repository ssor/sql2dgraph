package generator_v2

import (
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestGenerateCustomers(t *testing.T) {
    raw := []byte(`
insert  into customers( customerNumber , customerName , contactLastName , contactFirstName , phone , addressLine1 , addressLine2 , city , state , postalCode , country , salesRepEmployeeNumber ,creditLimit) values 
(103,'Atelier graphique','Schmitt','Carine','40.32.2555','54, rue Royale',NULL,'Nantes',NULL,'44000','France',1370,'21000.00'),
(112,'Signal Gift Stores','King','Jean','7025551838','8489 Strong St.',NULL,'Las Vegas','NV','83030','USA',1166,'71800.00');`)

    tableName, _, rows, err := ParseInsertSql(raw, nil)
    assert.Nil(t, err)
    customers, err := generateCustomers(tableName, rows)
    assert.Nil(t, err)
    assert.Equal(t, 2, len(customers))
    assert.Equal(t, 103, customers[0].CustomerNumber)
    assert.Equal(t, "Atelier graphique", customers[0].CustomerName)
    assert.Equal(t, "Schmitt", customers[0].ContactLastName)
    assert.Equal(t, "Carine", customers[0].ContactFirstName)
    assert.Equal(t, "40.32.2555", customers[0].Phone)
    assert.Equal(t, "54, rue Royale", customers[0].AddressLine1)
    assert.Equal(t, "", customers[0].AddressLine2)
    assert.Equal(t, "Nantes", customers[0].City)
    assert.Equal(t, "", customers[0].State)
    assert.Equal(t, "44000", customers[0].PostalCode)
    assert.Equal(t, "France", customers[0].Country)
    assert.Equal(t, 21000.0, customers[0].CreditLimit)
    assert.Equal(t, 1370, customers[0].Employee.EmployeeNumber)

    assert.Equal(t, 112, customers[1].CustomerNumber)
}

func TestDefaultCustomerValue(t *testing.T) {
    emptyCustomer := newEmptyCustomer()
    emptyCustomer.CustomerNumber = 0
    emptyCustomer.CreditLimit = 0
    raw, _ := json.Marshal(emptyCustomer)
    assert.Equal(t, 2, len(string(raw)))
}

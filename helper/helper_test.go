package helper

import (
    "context"
    "encoding/json"
    "github.com/dgraph-io/dgo"
    "github.com/dgraph-io/dgo/protos/api"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

var (
    dgClient *dgo.Dgraph
)

type Office struct {
    Uid          string `json:"uid,omitempty"`
    OfficeCode   string `json:"office_code,omitempty"`
    City         string `json:"city,omitempty"`
    Phone        string `json:"phone,omitempty"`
    AddressLine1 string `json:"address_line1,omitempty"`
    AddressLine2 string `json:"address_line2,omitempty"`
    State        string `json:"state,omitempty"`
    PostalCode   string `json:"postal_code,omitempty"`
    Country      string `json:"country,omitempty"`
    Territory    string `json:"territory,omitempty"`
}

func newOffice(code string) *Office {
    return &Office{
        OfficeCode: code,
    }
}

func (office *Office) QueryBy() []interface{} {
    return []interface{}{"office_code", office.OfficeCode}
}

func (office *Office) GetUidInfo() (index string, uid string) {
    uid = office.Uid
    return
}

func (office *Office) Schemes() string {
    var schemes Schemes
    schemes = schemes.Add(newSchemeStringExactIndex("office_code")).
        Add(newSchemeString("city")).
        Add(newSchemeString("phone"))

    return schemes.String()
}

func (office *Office) DependentObjectHasUid() bool {
    return true
}

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

func newEmployee(number int, reportsTo *Employee) *Employee {
    return &Employee{
        EmployeeNumber: number,
        ReportsTo:      reportsTo,
    }
}

func (employee *Employee) DependentObjectHasUid() bool {
    if employee.ReportsTo != nil {
        if len(employee.ReportsTo.Uid) <= 0 {
            return false
        }
    }
    return true
}

func (employee *Employee) GetUidInfo() (index string, uid string) {
    uid = employee.Uid
    return
}

// model without linked object
func TestAddAndUpdateOffice(t *testing.T) {
    office := newOffice("office-1")
    office.City = "beijing"
    uid, err := MutationObj(office, dgClient)
    assert.Nil(t, err)
    office.Uid = uid

    _, err = MutationObj(office, dgClient)
    assert.Nil(t, err)

    q := `query MyOffice($id: string){
        offices(func: uid($id)) {
            office_code
            city
        }
    }`
    resJson, err := QueryObjWithVars(q, map[string]string{"$id": office.Uid}, dgClient)
    assert.Nil(t, err)

    type Offices struct {
        MyOffice []*Office `json:"offices"`
    }
    offices := Offices{}
    err = json.Unmarshal(resJson, &offices)
    assert.Nil(t, err)
    assert.Equal(t, len(offices.MyOffice), 1, "should only one office")

    office.City = "nanjing"
    err = UpdateObj(office, dgClient)
    assert.Nil(t, err)

    resJson, err = QueryObjWithVars(q, map[string]string{"$id": office.Uid}, dgClient)
    assert.Equal(t, err, nil, "query failed")
    offices = Offices{}
    err = json.Unmarshal(resJson, &offices)
    assert.Nil(t, err)
    assert.Equal(t, len(offices.MyOffice), 1, "should only one office")
    assert.Equal(t, "nanjing", offices.MyOffice[0].City, "should be new city nanjing now, dump: "+string(resJson))
}

func TestRecursiveMutation(t *testing.T) {
    employee1 := newEmployee(1, nil)
    employee1.FirstName = "employee1"
    employee2 := newEmployee(2, employee1)
    employee2.FirstName = "employee2"

    uidOfEmployee1, err := MutationObj(employee1, dgClient)
    assert.Nil(t, err)
    employee1.Uid = uidOfEmployee1
    uidOfEmployee2, err := MutationObj(employee2, dgClient)
    assert.Nil(t, err)
    employee2.Uid = uidOfEmployee2

    q := `{
               employees(func: has(employee_number), orderasc: employee_number) {
                   employee_number
                   first_name
               }
           }`

    resJson, err := QueryObj(q, dgClient)
    assert.Nil(t, err)
    myEmployee := struct {
        Employees []*Employee `json:"employees"`
    }{}
    err = json.Unmarshal(resJson, &myEmployee)
    assert.Nil(t, err, string(resJson))
    assert.Equal(t, 2, len(myEmployee.Employees), string(resJson))
    assert.Equal(t, 1, myEmployee.Employees[0].EmployeeNumber, string(resJson))
    assert.Equal(t, 2, myEmployee.Employees[1].EmployeeNumber, string(resJson))
}

func TestMutateRaw(t *testing.T) {
    pb := []byte(`
     {
        "uid": "0xc244",
        "employee_number": 1003,
        "last_name": "Murphy",
        "first_name": "Diane",
        "extension": "x5800",
        "email": "dmurphy@classicmodelcars.com",
        "office_work_in": {
          "office_code": "1"
        },
        "job_title": "President"
      }
`)
    mu := &api.Mutation{
        CommitNow: true,
    }
    mu.SetJson = pb
    ctx := context.Background()
    assigned, err := dgClient.NewTxn().Mutate(ctx, mu)
    assert.Nil(t, err)
    assert.Equal(t, 1, len(assigned.Uids), assigned.String())

    ctx = context.Background()
    assigned, err = dgClient.NewTxn().Mutate(ctx, mu)
    assert.Nil(t, err)
    assert.Equal(t, 1, len(assigned.Uids), assigned.String())
}

func TestQueryUid(t *testing.T) {
    office := newOffice("office-3")
    office.City = "beijing"
    uid, err := MutationObj(office, dgClient)
    assert.Nil(t, err)

    err = Alter(office.Schemes(), dgClient)
    assert.Nil(t, err)

    queriedUid, err := QueryUid(office, dgClient)
    assert.Nil(t, err)
    assert.Equal(t, uid, queriedUid)

    office2 := newOffice("office-4")
    queriedUid, err = QueryUid(office2, dgClient)
    assert.Nil(t, err)
    assert.Equal(t, "", queriedUid)
}

// *** Drop relative data in dgraph, or test results may be affected ***
func TestMain(m *testing.M) {
    client, err := CreateDgClient("127.0.0.1:9080")
    if err != nil {
        panic(err)
    }
    dgClient = client
    os.Exit(m.Run())
}

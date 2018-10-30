package generator_v2

import (
    "github.com/stretchr/testify/assert"
    "os"
    "path"
    "testing"
)

func TestMutationObjs(t *testing.T) {
    err := os.Remove("uid.db")
    assert.Nil(t, err)

    err = DropDB()
    assert.Nil(t, err)

    err = AlterSchemas("employees", "customers")
    assert.Nil(t, err)

    err = MutationObjs(path.Join("test_data", "employees.sql"))
    assert.Nil(t, err)
    err = MutationObjs(path.Join("test_data", "customers.sql"))
    assert.Nil(t, err)
}

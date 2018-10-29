package generator_v2

import (
    "github.com/stretchr/testify/assert"
    "path"
    "testing"
)

func TestMutationObjs(t *testing.T) {
    err := MutationObjs(path.Join("test_data", "employees.sql"))
    assert.Nil(t, err)
    err = MutationObjs(path.Join("test_data", "customers.sql"))
    assert.Nil(t, err)

    err = AlterSchemas("employees", "customers")
    assert.Nil(t, err)
}

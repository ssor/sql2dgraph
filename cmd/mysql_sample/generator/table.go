package generator

import (
    "fmt"
    "github.com/xwb1989/sqlparser"
    "strings"
)

func ParseTables(src string) (Tables, error) {
    tablesRaw := strings.Split(src, "\n\r")
    var tables Tables
    for _, table := range tablesRaw {
        table = strings.Replace(table, "\n", " ", -1)
        table = strings.Replace(table, "\r", " ", -1)
        table = strings.Trim(table, " ")
        if len(table) > 0 {
            tableName, columns, err := ParseTable(table)
            if err != nil {
                fmt.Printf("parse table [%s] trigger error: %s", table, err)
                return nil, err
            }
            tables = append(tables, &Table{
                TableName: tableName,
                Columns:   columns,
            })
        }
    }
    return tables, nil
}

type Table struct {
    TableName string
    Columns   []*sqlparser.ColumnDefinition
}
type Tables []*Table

func (ts Tables) Find(name string) *Table {
    for _, table := range ts {
        if table.TableName == name {
            return table
        }
    }
    panic("cannot find table " + name)
}

func (ts Tables) SummaryDataType() []string {
    types := make(map[string]string)
    for _, table := range ts {
        for _, column := range table.Columns {
            types[column.Type.Type] = ""
        }
    }
    var tps []string
    for k := range types {
        tps = append(tps, k)
    }
    return tps
}

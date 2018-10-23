package generator

import (
    "bytes"
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "github.com/ssor/zlog"
    "github.com/xwb1989/sqlparser"
    "io/ioutil"
    "strings"
)

func ParseInsert(f string, tables Tables) (*Insert, error) {
    raw, err := ioutil.ReadFile(f)
    if err != nil {
        return nil, fmt.Errorf("read file %s trigger error: %s", f, err)
    }
    tableName, columnNames, rows, err := Parse(string(raw))
    if err != nil {
        return nil, fmt.Errorf("parse file %s trigger error: %s", f, err)
    }
    table := tables.Find(tableName)
    values := extractValuesFromRows(rows, table)
    insert := NewInsert(tableName, columnNames, values)
    insert.generateGraphSet()
    return insert, nil
}

func extractValuesFromRows(rows sqlparser.Values, table *Table) [][]string {
    var values [][]string
    for _, row := range rows {
        value := extractValueFromRow(row, table)
        if value != nil {
            values = append(values, value)
        } else {
            panic("unknown value type")
        }
    }
    return values
}

func extractValueFromRow(row sqlparser.ValTuple, table *Table) []string {
    var values []string
    vt := sqlparser.ValTuple(row)
    for index, expr := range vt {
        switch value := expr.(type) {
        case *sqlparser.SQLVal:
            values = append(values, string(value.Val))
        case *sqlparser.NullVal:
            typeName := table.Columns[index].Type.Type
            tn := mapKeyType(typeName)
            values = append(values, defaultValue(tn))
        default:
            zlog.Highlight("unknown value type")
            spew.Dump(value)
            return nil
        }
    }
    return values
}

func NewInsert(name string, columns []string, values [][]string) *Insert {
    insert := &Insert{
        TableName:   name,
        ColumnNames: columns,
        Values:      values,
        Units:       []*InsertUnit{},
    }
    return insert
}

type Insert struct {
    TableName   string
    ColumnNames []string
    Values      [][]string
    Units       []*InsertUnit
}

type InsertUnit struct {
    Index int
    Uid   string
    Sets  []string
}

func newInsertUnit(index int) *InsertUnit {
    return &InsertUnit{
        Index: index,
    }
}

func (unit *InsertUnit) Name(prefix string) string {
    return fmt.Sprintf("%s%d", prefix, unit.Index)
}

func (unit *InsertUnit) SetsString() string {
    sets := "{set{\r\n "
    sets += strings.Join(unit.Sets, "\r\n")
    sets += "\r\n}}"
    return sets
}

func (insert *Insert) RangeUnitSets(start, end int) string {
    if start > end {
        return ""
    }
    if end > len(insert.Units) {
        end = len(insert.Units)
    }
    if start < 0 {
        start = 0
    }
    units := insert.Units[start:end]
    sets := "{set{\r\n "
    for _, unit := range units {
        sets += strings.Join(unit.Sets, "\r\n")
        sets += "\r\n"
    }
    sets += "\r\n}}"
    return sets
}

func (insert *Insert) UnitSets(index int) string {
    if index > len(insert.Units) {
        return ""
    }
    if index < 0 {
        return ""
    }
    unit := insert.Units[index]
    return unit.SetsString()
}

func (insert *Insert) generateGraphSet() {
    for rowIndex, row := range insert.Values {
        var sets []string
        unit := newInsertUnit(rowIndex)
        for index, value := range row {
            className := insert.ClassFormat(rowIndex)
            bs := []byte(value)
            bs = bytes.Replace(bs, []byte("\n"), []byte(" "), -1)
            bs = bytes.Replace(bs, []byte("\r"), []byte(" "), -1)
            bs = bytes.Replace(bs, []byte("\""), []byte(" "), -1)
            sets = append(sets, newSet(className, insert.ColumnNames[index], string(bs)))
        }
        unit.Sets = sets
        insert.Units = append(insert.Units, unit)
    }
}

func (insert *Insert) ClassFormat(rowIndex int) string {
    src := fmt.Sprintf("%s%d", insert.TableName, rowIndex)
    return src
}

func (insert *Insert) FindRowIndexByValue(keyIndex int, value string) int {
    for index, row := range insert.Values {
        if row[keyIndex] == value {
            return index
        }
    }
    return -1
}

func (insert *Insert) KeyIndex(key string) int {
    for index := range insert.ColumnNames {
        if insert.ColumnNames[index] == key {
            return index
        }
    }
    return -1
}

type Inserts []*Insert

func (is Inserts) Find(name string) *Insert {
    for _, insert := range is {
        if insert.TableName == name {
            return insert
        }
    }
    panic("cannot find insert values " + name)
}

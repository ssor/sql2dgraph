package generator

import (
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "github.com/xwb1989/sqlparser"
    "time"
)

var ()

func init() {
    spew.Config.Indent = "    "
}

func GenerateIndices(tables Tables) []string {
    indicesRaw := make(map[string]string)
    for _, table := range tables {
        for _, column := range table.Columns {
            indicesRaw[column.Name.String()] = column.Type.Type
        }
    }

    var indices []string
    for columnName, columnTypeDefination := range indicesRaw {
        var indexType string
        switch columnTypeDefination {
        case "int", "smallint":
            indexType = "int"
        case "varchar", "text":
            indexType = "string"
        case "decimal":
            indexType = "float"
        case "date":
            indexType = "dateTime"
        default:
            panic("not support type of column " + columnTypeDefination)
        }
        indices = append(indices, fmt.Sprintf("%s: %s .", columnName, indexType))
    }
    return indices
}

func mapKeyType(src string) string {
    var indexType string
    switch src {
    case "int", "smallint":
        indexType = "int"
    case "varchar", "text":
        indexType = "string"
    case "decimal":
        indexType = "float"
    case "date":
        indexType = "dateTime"
    default:
        panic("not support type of column " + src)
    }
    return indexType
}

func defaultValue(src string) string {
    v := ""
    switch src {
    case "int", "float":
        v = "0"
    case "string":
    case "dateTime":
        v = time.Now().Format(time.RFC3339)
    default:
        panic("not support type of column " + src)
    }
    return v
}

func Parse(sql string) (tableName string, columnNames []string, values sqlparser.Values, e error) {
    stmt, err := sqlparser.ParseStrictDDL(sql)
    if err != nil {
        fmt.Println("parse sql failed, ", err)
        e = err
        return
    }

    switch stmt := stmt.(type) {
    case *sqlparser.Insert:
        tableName = stmt.Table.Name.String()
        for _, column := range stmt.Columns {
            columnNames = append(columnNames, column.String())
        }
        values = stmt.Rows.(sqlparser.Values)
        //for _, vt := range vts {
        //    spew.Dump(vt)
        //}
        return
    default:
        fmt.Println("=========unknown type=========")
    }
    return
}

func ParseTable(sql string) (string, []*sqlparser.ColumnDefinition, error) {
    stmt, err := sqlparser.ParseStrictDDL(sql)
    if err != nil {
        return "", nil, err
    }

    switch stmt := stmt.(type) {
    case *sqlparser.DDL:
        //fmt.Println("==========DDL===========")
        //spew.Dump(stmt.NewName.Name)
        return stmt.NewName.Name.String(), stmt.TableSpec.Columns, nil
    default:
        fmt.Println("=========unknown type=========")
        spew.Dump(stmt)
    }
    return "", nil, nil
}

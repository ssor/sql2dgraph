package generator_v2

import (
    "fmt"
    "github.com/xwb1989/sqlparser"
)

func ParseInsertSql(sql []byte, err error) (tableName string, columnNames []string, values sqlparser.Values, e error) {
    if err != nil {
        e = err
        return
    }

    stmt, err := sqlparser.ParseStrictDDL(string(sql))
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
        return
    default:
        logger.Warn("=========unknown type=========")
    }
    return
}

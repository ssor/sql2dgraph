package generator

import "fmt"

func newSet(className, relationName, value string) string {
    return fmt.Sprintf("_:%s <%s> \"%s\" .", className, relationName, value)
}

func newClassRelationSet(uid1, relationName, uid2 string) string {
    return fmt.Sprintf("<%s> <%s> <%s> .", uid1, relationName, uid2)
}

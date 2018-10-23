package generator

import (
    "fmt"
    "log"
)

type ClassRelation struct {
    From     string
    FromKey  string `yaml:"from_key"`
    Dest     string
    DestKey  string `yaml:"dest_key"`
    Relation string
}

func GenerateClassRelationSets(relations []*ClassRelation, inserts Inserts) []string {
    var sets []string
    for _, relation := range relations {
        set := GenerateClassRelationSet(relation, inserts)
        if set != nil {
            sets = append(sets, set...)
        } else {
            log.Println(fmt.Sprintf("******** relation %s <%s> %s has no sets", relation.From, relation.Dest, relation.Relation))
        }
    }
    return sets
}

func GenerateClassRelationSet(relation *ClassRelation, inserts Inserts) []string {
    tableFrom := inserts.Find(relation.From)
    keyIndexOfTableFrom := tableFrom.KeyIndex(relation.FromKey)
    if keyIndexOfTableFrom < 0 {
        return nil
    }

    tableDest := inserts.Find(relation.Dest)
    keyIndexOfTableDest := tableDest.KeyIndex(relation.DestKey)
    if keyIndexOfTableDest < 0 {
        return nil
    }

    var set []string
    for rowIndex, row := range tableFrom.Values {
        v := row[keyIndexOfTableFrom]
        destRowIndex := tableDest.FindRowIndexByValue(keyIndexOfTableDest, v)
        if destRowIndex < 0 {
            continue
        }
        unitSrc := tableFrom.Units[rowIndex]
        unitDest := tableDest.Units[destRowIndex]
        s := newClassRelationSet(unitSrc.Uid, relation.Relation, unitDest.Uid)
        set = append(set, s)
    }
    return set
}

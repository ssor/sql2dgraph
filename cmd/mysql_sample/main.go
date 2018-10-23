package main

import (
    "bytes"
    "fmt"
    "github.com/mkideal/cli"
    "github.com/ssor/sql2graphql/cmd/mysql_sample/generator"
    "github.com/ssor/zlog"
    "github.com/tidwall/gjson"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "net/http"
    "os"
    "path"
    "path/filepath"
    "strconv"
    "strings"
)

/*
    1. mutate tables and get uid, write this map to file
    2. create relations supported by this map
    3. mutate relations
    4. add scheme
*/

var (
    inserts   generator.Inserts
    tables    generator.Tables
    relations []*generator.ClassRelation
)

type Args struct {
    cli.Helper
    Command   int    `cli:"cmd"  usage:"what to do? 1: mutation obj; 2: mutation relations; 3: add scheme"  dft:"0"`
    OutputDir string `cli:"output-dir"  usage:"dir for data output, default is data/output"  dft:"data/output"`
    RootDir   string `cli:"root"  usage:"the root dir of tables and table values, default is data/mysql_sample"  dft:"data/mysql_sample"`
    //FileOfTables string `cli:"table"  usage:"path of tables, default is data/mysql_sample/tables/tables.sql"  dft:"data/mysql_sample/tables/tables.sql"`
}

func main() {
    cli.Run(new(Args), func(ctx *cli.Context) error {
        args := ctx.Argv().(*Args)
        err := os.MkdirAll(args.OutputDir, os.ModePerm)
        if err != nil {
            panic(err)
        }
        zlog.Info("cmd -> ", args.Command)

        switch args.Command {
        case 1:
            mutationObjs(args.RootDir, args.OutputDir)
        case 2:
            mutateRelations(args.RootDir, args.OutputDir)
        case 3:
            alterSchemes(args.RootDir, args.OutputDir)
        default:
        }
        return nil
    })
}

func getTablePath(rootDir string) string {
    tablePath := path.Join(rootDir, "tables", "tables.sql")
    return tablePath
}

func alterSchemes(rootDir, outputDir string) {
    tables = loadTables(rootDir)
    zlog.Info("load tables OK")
    indices := generator.GenerateIndices(tables)
    indicesSet := outputIndices(indices, outputDir)
    err := alterIndices(indicesSet)
    if err != nil {
        panic(err)
    }
}

func mutateRelations(rootDir, outputDir string) {
    tables = loadTables(rootDir)
    zlog.Info("load tables OK")

    inserts, err := loadInsertFiles(tables, rootDir)
    if err != nil {
        panic(err)
    }
    validateData()
    zlog.Info("data validate OK")

    err = loadUidInfo(tables, inserts, outputDir)
    if err != nil {
        panic(err)
    }
    zlog.Info("load UID OK")

    relations = parseRelations(rootDir)
    for _, relation := range relations {
        sets := generator.GenerateClassRelationSet(relation, inserts)
        err := outputClassRelations(relation.From, sets, outputDir)
        if err != nil {
            panic(err)
        }
    }

    zlog.Info("output relation OK")
}

func mutationObjs(rootDir, outputDir string) {
    tables = loadTables(rootDir)
    zlog.Info("load tables OK")
    inserts, err := loadInsertFiles(tables, rootDir)
    if err != nil {
        panic(err)
    }
    zlog.Info("load insert values OK")
    validateData()
    zlog.Info("data validate OK")

    for _, insert := range inserts {
        outputSets(insert, outputDir)
    }
    zlog.Info("output obj sets OK")

    for _, objInsert := range inserts {
        for _, unit := range objInsert.Units {
            uid, err := mutation(unit.Name(objInsert.TableName), unit.SetsString())
            if err != nil {
                fmt.Println(strings.Repeat("-", 64))
                fmt.Println(unit.SetsString())
                fmt.Println(strings.Repeat("-", 64))
                panic(err)
            }
            unit.Uid = uid
        }
    }

    zlog.Info("mutate obj OK")
    for _, insert := range inserts {
        outputSetsUid(insert, outputDir)
    }
    zlog.Info("output obj uid OK")
}

func mutateSets(sets string) (string, error) {
    client := &http.Client{}
    req, err := http.NewRequest("POST", "http://localhost:8080/mutate", strings.NewReader(sets))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("X-Dgraph-CommitNow", "true")
    resp, err := client.Do(req)
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    errors := gjson.Get(string(body), "errors")
    if len(errors.Array()) > 0 {
        return "", fmt.Errorf("failed mutate for %s", errors.String())
    }

    return string(body), nil
}

func mutation(name, sets string) (uid string, e error) {
    body, err := mutateSets(sets)
    if err != nil {
        e = err
        return
    }
    value := gjson.Get(body, "data.uids")

    value.ForEach(func(key, value gjson.Result) bool {
        if key.String() == name {
            uid = value.String()
        }
        return true
    })

    if len(uid) <= 0 {
        printLine()
        zlog.PrettyJson([]byte(body))
        zlog.Info(name)
        zlog.Info(sets)
        printLine()
        //fmt.Println(value.String())
    }
    return
}

func printLine() {
    println(strings.Repeat("-", 100))
}

func parseRelations(rootDir string) []*generator.ClassRelation {
    fpath := path.Join(rootDir, "tables", "relations.yaml")
    relationsRaw, err := ioutil.ReadFile(fpath)
    if err != nil {
        panic(err)
    }
    var relations []*generator.ClassRelation
    err = yaml.Unmarshal(relationsRaw, &relations)
    if err != nil {
        panic(err)
    }
    //spew.Dump(relations)
    return relations
}

func outputClassRelations(prefix string, allSets []string, outputDir string) error {
    splittedSets := splitSetN(allSets, 500)
    for index, sets := range splittedSets {
        buffer := bytes.NewBuffer([]byte{})
        buffer.WriteString("{set{\n")
        for index := range sets {
            buffer.WriteString(sets[index])
            buffer.WriteString("\n")
        }
        buffer.WriteString("}}")
        err := ioutil.WriteFile(filepath.Join(outputDir, fmt.Sprintf("relations_%s_%d.json", prefix, index)), buffer.Bytes(), os.ModePerm)
        if err != nil {
            panic(err)
        }
        zlog.Infof("output %s OK", prefix)

        _, err = mutateSets(buffer.String())
        if err != nil {
            zlog.Errorf("mutate relations failed: %s", err)
            printLine()
            zlog.Info(allSets)
            printLine()
            return err
        }

        zlog.Infof("mutate %s OK", prefix)
    }
    return nil
}

func splitSetN(sets []string, maxRow int) [][]string {
    var nRows [][]string
    totalLength := len(sets)
    if totalLength <= maxRow {
        nRows = append(nRows, sets)
        return nRows
    }
    count := 0
    i := 0
    for {
        if (i + maxRow) > totalLength {
            break
        }
        slice := sets[i : i+maxRow]
        count += len(slice)
        nRows = append(nRows, slice)
        i += maxRow
    }

    slice := sets[i:]
    nRows = append(nRows, slice)
    count += len(slice)
    if count != totalLength {
        panic(fmt.Sprintf("splitSet error total %d != %d", totalLength, count))
    }
    return nRows
}

func outputIndices(indices []string, outputDir string) string {
    buffer := bytes.NewBuffer([]byte{})
    for index := range indices {
        buffer.WriteString(indices[index])
        buffer.WriteString("\n")
    }
    err := ioutil.WriteFile(filepath.Join(outputDir, "indices.json"), buffer.Bytes(), os.ModePerm)
    if err != nil {
        panic(err)
    }
    zlog.Info("output Indices OK")
    return buffer.String()
}

func alterIndices(indices string) error {
    body, err := alter(indices)
    if err != nil {
        zlog.Errorf("alter schemes failed: %s", err)
        printLine()
        zlog.Info(body)
        zlog.Info(indices)
        printLine()
        return err
    }
    zlog.Info("alter indices OK")
    return nil
}

func alter(sets string) (string, error) {
    client := &http.Client{}
    req, err := http.NewRequest("POST", "http://localhost:8080/alter", strings.NewReader(sets))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("X-Dgraph-CommitNow", "true")
    resp, err := client.Do(req)
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    errors := gjson.Get(string(body), "errors")
    if len(errors.Array()) > 0 {
        return "", fmt.Errorf("failed mutate for %s", errors.String())
    }
    return string(body), nil
}

func outputSetsUid(insert *generator.Insert, outputDir string) {
    buffer := bytes.NewBuffer([]byte{})
    units := insert.Units
    for _, unit := range units {
        if len(unit.Uid) > 0 {
            buffer.WriteString(fmt.Sprintf("%d:%s", unit.Index, unit.Uid))
            buffer.WriteString(";")
        }
    }
    err := ioutil.WriteFile(filepath.Join(outputDir, insert.TableName+"_uid.json"), buffer.Bytes(), os.ModePerm)
    if err != nil {
        panic(err)
    }
}

func outputSets(insert *generator.Insert, outputDir string) {
    prefix := insert.TableName
    units := insert.Units
    buffer := bytes.NewBuffer([]byte{})
    buffer.WriteString("{set{\n")
    for _, unit := range units {
        for _, set := range unit.Sets {
            buffer.WriteString(set)
            buffer.WriteString("\r\n")
        }
    }
    buffer.WriteString("}}")
    err := ioutil.WriteFile(filepath.Join(outputDir, prefix+".json"), buffer.Bytes(), os.ModePerm)
    if err != nil {
        panic(err)
    }
}

func validateData() {
    for _, insert := range inserts {
        tbName := insert.TableName

        exists := false
        for _, tb := range tables {
            if tb.TableName == tbName {
                exists = true
                break
            }
        }
        if exists == false {
            panic(tbName + "no found")
        }
    }
}

func loadUidInfo(tables generator.Tables, inserts generator.Inserts, outputDir string) (error) {
    for _, table := range tables {
        fpath := filepath.Join(outputDir, table.TableName+"_uid.json")
        bs, err := ioutil.ReadFile(fpath)
        if err != nil {
            return err
        }
        maps := strings.Split(string(bs), ";")
        for index := range maps {
            raw := maps[index]
            if len(raw) <= 0 {
                continue
            }
            indexUid := strings.Split(raw, ":")
            insert := inserts.Find(table.TableName)
            rowIndex, _ := strconv.Atoi(indexUid[0])
            unit := insert.Units[rowIndex]
            unit.Uid = indexUid[1]
        }
    }
    return nil
}

func loadInsertFiles(tables generator.Tables, rootPath string) (generator.Inserts, error) {
    var inserts generator.Inserts
    for _, table := range tables {
        fpath := filepath.Join(rootPath, "values", table.TableName+".sql")
        insert, err := generator.ParseInsert(fpath, tables)
        if err != nil {
            return nil, fmt.Errorf("parse file %s trigger error: %s", fpath, err)
        }
        inserts = append(inserts, insert)
    }

    return inserts, nil
}

func loadTables(rootDir string) generator.Tables {
    if tables != nil {
        return tables
    }

    tablesPath := getTablePath(rootDir)
    raw, err := ioutil.ReadFile(tablesPath)
    if err != nil {
        fmt.Printf("read file %s trigger error: %s \n", tablesPath, err)
        return nil
    }
    tables, err := generator.ParseTables(string(raw))
    if err != nil {
        fmt.Printf("parse tables error: %s \n", err)
        return nil
    }
    return tables
}

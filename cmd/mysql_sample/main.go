package main

import (
    "github.com/mkideal/cli"
    "github.com/ssor/sql2graphql/cmd/mysql_sample/generator"
    "github.com/ssor/zlog"
    "os"
)

/*
    1. mutate tables and get uid, write this map to file
    2. create relations supported by this map
    3. mutate relations
    4. add scheme
*/

//var (
//    inserts   generator.Inserts
//    tables    generator.Tables
//    relations []*generator.ClassRelation
//)

type Args struct {
    cli.Helper
    Command   int    `cli:"cmd"  usage:"what to do? 1: mutation obj; 2: mutation relations; 3: add scheme"  dft:"0"`
    OutputDir string `cli:"output-dir"  usage:"dir for data output, default is data/output"  dft:"data/output"`
    RootDir   string `cli:"root"  usage:"the root dir of tables and table values, default is data/mysql_sample"  dft:"data/mysql_sample"`
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
            generator.MutationObjs(args.RootDir, args.OutputDir)
        case 2:
            generator.MutateRelations(args.RootDir, args.OutputDir)
        case 3:
            generator.AlterSchemes(args.RootDir, args.OutputDir)
        default:
        }
        return nil
    })
}

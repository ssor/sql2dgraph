package main

import (
    "github.com/mkideal/cli"
    "github.com/ssor/sql2graphql/cmd/mysql_sample/generator"
    "github.com/ssor/sql2graphql/cmd/mysql_sample/generator.v2"
    "github.com/ssor/zlog"
    "os"
    "path"
)

type Args struct {
    cli.Helper
    Version   int    `cli:"v,version"  usage:"generator version, version 1 will output rdf sets, version 2 will not"  dft:"1"`
    Command   int    `cli:"cmd"  usage:"what to do? 1: mutation obj; 2: mutation relations; 3: add scheme"  dft:"0"`
    OutputDir string `cli:"output-dir"  usage:"dir for data output, default is cmd/mysql_sample/data/output"  dft:"cmd/mysql_sample/data/output"`
    RootDir   string `cli:"root"  usage:"the root dir of tables and table values, default is cmd/mysql_sample/data"  dft:"cmd/mysql_sample/data"`
}

func main() {
    cli.Run(new(Args), func(ctx *cli.Context) error {
        args := ctx.Argv().(*Args)
        zlog.WithFields(map[string]interface{}{
            "version":    args.Version,
            "command":    args.Command,
            "output_dir": args.OutputDir,
            "root_dir":   args.RootDir,
        }).Info("Run with args:")

        switch args.Version {
        case 1:
            v1Handler(args)
        case 2:
            v2Handler(args)
        default:
            zlog.Failed("unknown version, nothing done")
        }
        return nil
    })
}

func v2Handler(args *Args) {
    err := generator_v2.MutationObjs(path.Join(args.RootDir, "values", "employees.sql"), path.Join(args.RootDir, "values", "customers.sql"))
    if err != nil {
        zlog.Error(err)
    }
}

func v1Handler(args *Args) {
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
}

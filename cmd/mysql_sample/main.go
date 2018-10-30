package main

import (
    "github.com/mkideal/cli"
    "github.com/ssor/sql2graphql/cmd/mysql_sample/generator.v2"
    "github.com/ssor/zlog"
    "path"
)

type Args struct {
    cli.Helper
    Command int    `cli:"cmd"  usage:"what to do? 0: import data; 1: drop db"  dft:"-1"`
    RootDir string `cli:"root"  usage:"the root dir of tables and table values, default is cmd/mysql_sample/data"  dft:"cmd/mysql_sample/data"`
    //Version   int    `cli:"v,version"  usage:"generator version, version 1 will output rdf sets, version 2 will not"  dft:"2"`
    //OutputDir string `cli:"output-dir"  usage:"dir for data output, default is cmd/mysql_sample/data/output"  dft:"cmd/mysql_sample/data/output"`
}

func main() {
    cli.Run(new(Args), func(ctx *cli.Context) error {
        args := ctx.Argv().(*Args)
        zlog.WithFields(map[string]interface{}{
            "command":  args.Command,
            "root_dir": args.RootDir,
            //"output_dir": args.OutputDir,
            //"version":    args.Version,
        }).Info("Run with args:")

        switch args.Command {
        case 0:
            //v1Handler(args)
            v2Handler(args)
        case 1:
            generator_v2.DropDB()
        default:
            zlog.Failed("unknown version, nothing done")
        }
        return nil
    })
}

func v2Handler(args *Args) {
    err := generator_v2.AlterSchemas("employees", "customers")
    if err != nil {
        zlog.Error(err)
    }

    err = generator_v2.MutationObjs(path.Join(args.RootDir, "values", "employees.sql"), path.Join(args.RootDir, "values", "customers.sql"))
    if err != nil {
        zlog.Error(err)
    }
}

//
//func v1Handler(args *Args) {
//    err := os.MkdirAll(args.OutputDir, os.ModePerm)
//    if err != nil {
//        panic(err)
//    }
//    zlog.Info("cmd -> ", args.Command)
//
//    switch args.Command {
//    case 1:
//        generator.MutationObjs(args.RootDir, args.OutputDir)
//    case 2:
//        generator.MutateRelations(args.RootDir, args.OutputDir)
//    case 3:
//        generator.AlterSchemes(args.RootDir, args.OutputDir)
//    default:
//    }
//}

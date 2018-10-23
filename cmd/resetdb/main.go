package main

import (
    "context"
    "fmt"
    "github.com/dgraph-io/dgo"
    "github.com/dgraph-io/dgo/protos/api"
    "github.com/mkideal/cli"
    "github.com/ssor/zlog"
    "google.golang.org/grpc"
)

type Args struct {
    cli.Helper
    Host string `cli:"host"  usage:"dgraph host"  dft:"127.0.0.1"`
    Port int    `cli:"port"  usage:"Default 9080"  dft:"9080"`
}

func main() {
    cli.Run(new(Args), func(ctx *cli.Context) error {
        args := ctx.Argv().(*Args)
        conn, err := grpc.Dial(fmt.Sprintf("%s:%d", args.Host, args.Port), grpc.WithInsecure())
        if err != nil {
            zlog.Failed("While trying to dial gRPC")
            return err
        }
        defer conn.Close()

        dc := api.NewDgraphClient(conn)
        dg := dgo.NewDgraphClient(dc)

        op := api.Operation{
            DropAll: true,
        }
        dgctx := context.Background()
        if err := dg.Alter(dgctx, &op); err != nil {
            zlog.Failed(err)
            return err
        }

        zlog.Success("dgraph drop success")
        return nil
    })
}

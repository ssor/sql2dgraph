package server

import (
    "github.com/dgraph-io/dgo"
    "github.com/ssor/sql2graphql/helper"
    "github.com/ssor/zlog"
)

var (
    logger = zlog.New("chain", "stream")

    dgClient *dgo.Dgraph
)

func init() {
    logger.Info("init ...")
    client, err := helper.CreateDgClient("127.0.0.1:9080")
    if err != nil {
        panic(err)
    }
    dgClient = client
}

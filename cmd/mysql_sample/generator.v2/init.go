package generator_v2

import (
    "github.com/dgraph-io/dgo"
    "github.com/ssor/sql2graphql/helper"
    "github.com/ssor/zlog"
)

var (
    logger = zlog.New("mysql_sample", "generator.v2")

    dgClient         *dgo.Dgraph
    //kvdb             *bolt.DB
    //normalBucketName = []byte("kv_bucket")
)

func init() {
    logger.Info("init ...")
    client, err := helper.CreateDgClient("127.0.0.1:9080")
    if err != nil {
        panic(err)
    }
    dgClient = client
    //
    //db, err := bolt.Open("uid.db", 0600, nil)
    //if err != nil {
    //    panic("init db failed for " + err.Error())
    //}
    //kvdb = db
}

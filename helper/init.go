package helper

import (
    "github.com/boltdb/bolt"
    "github.com/ssor/zlog"
)

var (
    logger = zlog.New("graphql", "helper")

    kvdb             *bolt.DB
    normalBucketName = []byte("kv_bucket")
)

func init() {
    db, err := bolt.Open("uid.db", 0600, nil)
    if err != nil {
        panic("init db failed for " + err.Error())
    }
    kvdb = db
}

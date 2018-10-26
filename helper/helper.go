package helper

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "github.com/dgraph-io/dgo"
    "github.com/dgraph-io/dgo/protos/api"
    "github.com/tidwall/gjson"
)

var (
    ErrNil      = fmt.Errorf("nil")
    ErrUidUnset = fmt.Errorf("uid of dependent object not set")
    //ErrUnmutatable=fmt.Errorf("Object has no UID")
)

type Mutatable interface {
    DependentObjectHasUid() bool
    GetUidInfo() (index string, uid string)
}

func MutationObj(obj Mutatable, client *dgo.Dgraph) (uid string, e error) {
    if obj.DependentObjectHasUid() == false {
        e = ErrUidUnset
        return
    }

    mu := &api.Mutation{
        CommitNow: true,
    }

    pb, err := json.Marshal(obj)
    if err != nil {
        logger.Warn("marshal obj failed: ", err)
        e = err
        return
    }

    //logger.Highlight("mutate ->")
    //zlog.PrettyJson(pb)
    mu.SetJson = pb
    ctx := context.Background()
    assigned, err := client.NewTxn().Mutate(ctx, mu)
    if err != nil {
        logger.Warn("While trying to mutate failed: ", err)
        e = err
        return
    }

    _, existedUid := obj.GetUidInfo()
    if len(existedUid) > 0 {
        uid = existedUid
    } else {
        newUid, exists := assigned.Uids["blank-0"]
        if exists == false {
            spew.Dump(assigned)
            e = fmt.Errorf("uid not returned in mutation")
            return
        }
        uid = newUid
    }
    return
}

func QueryObj(query string, client *dgo.Dgraph) (json []byte, e error) {
    ctx := context.Background()
    resp, err := client.NewTxn().Query(ctx, query)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        e = err
        return
    }
    json = resp.Json
    return
}

func QueryObjWithVars(query string, variables map[string]string, client *dgo.Dgraph) (json []byte, e error) {
    ctx := context.Background()
    resp, err := client.NewTxn().QueryWithVars(ctx, query, variables)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        e = err
        return
    }
    json = resp.Json
    return
}

func UpdateObj(obj interface{}, client *dgo.Dgraph) (e error) {
    mu := &api.Mutation{
        CommitNow: true,
    }
    pb, err := json.Marshal(obj)
    if err != nil {
        logger.Warn("marshal obj failed: ", err)
        e = err
        return
    }

    mu.SetJson = pb
    ctx := context.Background()
    assigned, err := client.NewTxn().Mutate(ctx, mu)
    if err != nil {
        logger.Warn("While trying to update failed: ", err)
        spew.Dump(assigned)
        e = err
        return
    }
    return
}

type Queryable interface {
    QueryBy() []interface{}
    Schemes() string
}

func QueryUid(qb Queryable, client *dgo.Dgraph) (uid string, e error) {
    query := fmt.Sprintf(`{
        query_uid(func: eq(%s, "%s")) {
            uid
        }
    }`, qb.QueryBy()...)
    ctx := context.Background()
    resp, err := client.NewTxn().Query(ctx, query)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        spew.Dump(query)
        e = err
        return
    }

    uid = gjson.Get(string(resp.Json), "query_uid.0.uid").String()
    return
}

func Alter(schemes string, client *dgo.Dgraph) (e error) {
    op := &api.Operation{}
    op.Schema = schemes
    ctx := context.Background()
    err := client.Alter(ctx, op)
    if err != nil {
        logger.Warn("While try to alter scheme failed: ", err)
        e = err
        return
    }
    return
}

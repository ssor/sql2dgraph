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
    ErrNeedUid  = fmt.Errorf("update object need UID")
)

type Mutatable interface {
    DependentObjectHasUid() bool
    GetUidInfo() (index string, uid string)
    SetUid(uid string)
    QueryBy() (name string, value string)
}

func MutationObj(obj Mutatable, client *dgo.Dgraph) (uid string, e error) {
    if obj.DependentObjectHasUid() == false {
        e = ErrUidUnset
        return
    }
    // If uid can be found in local db, we will query use it
    uidExisted, err := GetUidByIndex(obj, client)
    if err != nil {
        e = err
        return
    }

    obj.SetUid(uidExisted)

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
        logger.Warn("While trying to mutate failed: ", err)
        e = err
        return
    }

    index, existedUid := obj.GetUidInfo()
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

        saveIndexUidMapToLocalDb(index, uid) // this may have err, but we can ignore it as the uid has been stored in dgraph
    }
    return
}

// UpdateObj update data of node with uid
// Uid will be checked first, and if there is no uid, update will failed
func UpdateObj(obj Mutatable, client *dgo.Dgraph) (e error) {
    _, uid := obj.GetUidInfo()
    if len(uid) <= 0 {
        // uid not set, we should get this uid from local db
        // If uid cannot found in local db, we will query it from dgraph
        uid, err := GetUidByIndex(obj, client)
        if err != nil {
            logger.Warn("get uid for object failed: ", err)
            return ErrNeedUid
        }
        obj.SetUid(uid)
    }

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

// QueryObj do query from dgraph
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

// QueryObjWithVars do query from dgraph with paras
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

type Queryable interface {
    QueryBy() (string, string)
}

func QueryUid(qb Queryable, client *dgo.Dgraph) (uid string, e error) {
    name, value := qb.QueryBy()
    query := fmt.Sprintf(`{
        query_uid(func: eq(%s, "%s")) {
            uid
        }
    }`, name, value)
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

func Alter(schemes Schemas, client *dgo.Dgraph) (e error) {
    op := &api.Operation{}
    op.Schema = schemes.String()
    ctx := context.Background()
    err := client.Alter(ctx, op)
    if err != nil {
        logger.Warn("While try to alter scheme failed: ", err)
        e = err
        return
    }
    return
}

func DropDB(client *dgo.Dgraph) error {
    op := api.Operation{
        DropAll: true,
    }
    ctx := context.Background()
    if err := client.Alter(ctx, &op); err != nil {
        logger.Failed(err)
        return err
    }

    logger.Success("dgraph drop success")
    return nil
}

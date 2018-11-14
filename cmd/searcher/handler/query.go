package handler

import (
    "context"
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "github.com/dgraph-io/dgo"
    "github.com/tidwall/gjson"
)

//func QueryTxCount(ctx *gin.Context) {
//}

type BlockInfo struct {
    HashID string `json:"hash_id"`
    Height int    `json:"height"`
}

type Transaction struct {
    HashID       string     `json:"hash_id"`
    IndexInBlock int        `json:"index_in_block"`
    BlockHeight  int        `json:"block_height"`
    Block        *BlockInfo `json:"block_detail"`
}

type QueryByHashResult struct {
    Blocks []*BlockInfo   `json:"blocks"`
    Txs    []*Transaction `json:"txs"`
}

func queryAllBlocks(from int, client *dgo.Dgraph) ([]*BlockInfo, error) {
    query := fmt.Sprintf(`
        {
            query_all_blocks(func: has(block_hash_id), first:60, offset:%d) {
                block_hash_id
                block_height
            }
        }
    `, from)
    ctx := context.Background()
    resp, err := client.NewTxn().Query(ctx, query)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        spew.Dump(query)
        return nil, err
    }
    var blocks []*BlockInfo
    blocksResult := gjson.Get(string(resp.Json), "query_all_blocks")
    for _, block := range blocksResult.Array() {
        var bi BlockInfo
        bi.HashID = block.Get("block_hash_id").String()
        bi.Height = int(block.Get("block_height").Int())
        blocks = append(blocks, &bi)
    }
    return blocks, nil
}

func queryByHash(keyword string, client *dgo.Dgraph) (*QueryByHashResult, error) {
    if len(keyword) < 3 {
        return nil, fmt.Errorf("keyword is to simple")
    }
    query := fmt.Sprintf(`
    {
        query_block_by_hash(func: regexp(block_hash_id, /\S*%s\S*/)) {
                block_hash_id
                block_height
        }

        query_tx_by_hash(func: regexp(tx_hash_id, /\S*%s\S*/)) {
            tx_hash_id
            index_in_block
             tx_in_block : tx.block{
                block_hash_id
                block_height
            }
        }
    }`, keyword, keyword)
    logger.Info("queryByHash: ", query)
    ctx := context.Background()
    resp, err := client.NewTxn().Query(ctx, query)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        spew.Dump(query)
        return nil, err
    }
    var result QueryByHashResult
    blocksResult := gjson.Get(string(resp.Json), "query_block_by_hash")
    for _, block := range blocksResult.Array() {
        var bi BlockInfo
        bi.HashID = block.Get("block_hash_id").String()
        bi.Height = int(block.Get("block_height").Int())
        result.Blocks = append(result.Blocks, &bi)
    }

    txsResult := gjson.Get(string(resp.Json), "query_tx_by_hash")
    for _, block := range txsResult.Array() {
        var tx Transaction
        tx.HashID = block.Get("tx_hash_id").String()
        tx.IndexInBlock = int(block.Get("index_in_block").Int())
        var bi BlockInfo
        bi.HashID = block.Get("tx_in_block.0.block_hash_id").String()
        bi.Height = int(block.Get("tx_in_block.0.block_height").Int())
        tx.Block = &bi
        tx.BlockHeight = bi.Height
        result.Txs = append(result.Txs, &tx)
    }
    //spew.Dump(result)
    return &result, nil
}

func queryAllTxsInBlock(height int, client *dgo.Dgraph) ([]*Transaction, error) {
    query := fmt.Sprintf(`
    {
         query_all_txs_in_block(func: eq(tx.block_height, %d))  {
            tx_hash_id
            index_in_block
            height : tx.block_height
         }
    }
    `, height)
    ctx := context.Background()
    resp, err := client.NewTxn().Query(ctx, query)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        spew.Dump(query)
        return nil, err
    }
    var txs []*Transaction
    txsResult := gjson.Get(string(resp.Json), "query_all_txs_in_block")
    for _, txr := range txsResult.Array() {
        var tx Transaction
        tx.HashID = txr.Get("tx_hash_id").String()
        tx.IndexInBlock = int(txr.Get("index_in_block").Int())
        tx.BlockHeight = int(txr.Get("height").Int())
        txs = append(txs, &tx)
        //spew.Dump(txr)
    }

    return txs, nil
}

func queryTxCount(client *dgo.Dgraph) int {
    query := fmt.Sprintf(`
    {
        query_tx_count(func: has(tx_hash_id)) {
            count(uid)
        }
    }`)
    ctx := context.Background()
    resp, err := client.NewTxn().Query(ctx, query)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        spew.Dump(query)
        return 0
    }

    count := gjson.Get(string(resp.Json), "query_tx_count.0.count").Int()
    logger.Debugf("get tx count result: [%d]", count)
    return int(count)
}

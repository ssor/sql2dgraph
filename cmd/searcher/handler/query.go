package handler

import (
    "context"
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "github.com/dgraph-io/dgo"
    "github.com/tidwall/gjson"
)

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

type QueryResult struct {
    Blocks []*BlockInfo   `json:"blocks"`
    Txs    []*Transaction `json:"txs"`
}

//func queryTxRangeInBlock(blockHeight, from, to int) (*QueryResult, error) {
//    return &QueryResult{}, nil
//}

func queryBlockRange(from, to int, client *dgo.Dgraph) (*QueryResult, error) {
    query := fmt.Sprintf(`
        {
            query_block_by_height(func: has(block_height)) @filter(ge(block_height, %d) and le(block_height, %d)) {
                    block_hash_id
                    block_height
            }
        }
    `, from, to)
    ctx := context.Background()
    resp, err := client.NewTxn().Query(ctx, query)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        spew.Dump(query)
        return nil, err
    }

    var result QueryResult
    blocksResult := gjson.Get(string(resp.Json), "query_block_by_height")
    for _, block := range blocksResult.Array() {
        var bi BlockInfo
        bi.HashID = block.Get("block_hash_id").String()
        bi.Height = int(block.Get("block_height").Int())
        result.Blocks = append(result.Blocks, &bi)
    }

    return &result, nil
}

//
//func querySingleBlockByHeight(height int, client *dgo.Dgraph) (*QueryResult, error) {
//    query := fmt.Sprintf(`
//    {
//        query_block_by_height(func: eq(block_height, %d)) {
//                block_hash_id
//                block_height
//        }
//    }
//    `, height)
//    ctx := context.Background()
//    resp, err := client.NewTxn().Query(ctx, query)
//    if err != nil {
//        logger.Warn("While try to query failed: ", err)
//        spew.Dump(query)
//        return nil, err
//    }
//    var bi BlockInfo
//    firstBlockResult := gjson.Get(string(resp.Json), "query_block_by_height.0")
//    bi.HashID = firstBlockResult.Get("block_hash_id").String()
//    bi.Height = int(firstBlockResult.Get("block_height").Int())
//
//    var result QueryResult
//    result.Blocks = append(result.Blocks, &bi)
//
//    return &result, nil
//}

func queryAllBlocks(from int, client *dgo.Dgraph) (*QueryResult, error) {
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
    var result QueryResult
    result.Blocks = blocks
    return &result, nil
}

func queryByHash(keyword string, client *dgo.Dgraph) (*QueryResult, error) {
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
    ctx := context.Background()
    resp, err := client.NewTxn().Query(ctx, query)
    if err != nil {
        logger.Warn("While try to query failed: ", err)
        spew.Dump(query)
        return nil, err
    }
    var result QueryResult
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
    if result.Txs == nil && result.Blocks == nil {
        logger.Info("queryByHash: ", query)
    }
    //spew.Dump(result)
    return &result, nil
}

func mapFromtoToOffset(from, to int) (first, offset int) {
    offset = from - 1
    first = to - from + 1
    if offset < 0 {
        offset = 0
    }
    if first <= 0 {
        first = 1
    }
    if first > 60 {
        first = 60
    }
    return
}

func queryTxsInBlock(height, from, to int, client *dgo.Dgraph) (*QueryResult, error) {
    first, offset := mapFromtoToOffset(from, to)
    query := fmt.Sprintf(`
    {
         query_all_txs_in_block(func: eq(tx.block_height, %d), first:%d, offset:%d)  {
            tx_hash_id
            index_in_block
            height : tx.block_height
            tx_in_block : tx.block{
                block_hash_id
                block_height
            }
         }
    }
    `, height, first, offset)
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
        var bi BlockInfo
        bi.HashID = txr.Get("tx_in_block.0.block_hash_id").String()
        bi.Height = int(txr.Get("tx_in_block.0.block_height").Int())
        tx.Block = &bi
        txs = append(txs, &tx)
        //spew.Dump(txr)
    }
    if len(txs) <= 0 {
        logger.Info("no tx get for query:\n", query)
    }
    var result QueryResult
    result.Txs = txs
    return &result, nil
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
    if count <= 0 {
        logger.Debugf("get tx count result: [%d] from query \n%s", count, query)
    }
    return int(count)
}

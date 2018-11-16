package handler

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestQueryTxCount(t *testing.T) {
    count := queryTxCount(dgClient)
    assert.Equal(t, true, count > 1)
}

func TestQueryTxsInBlock(t *testing.T) {
    result, err := queryTxsInBlock(2, 0, 100, dgClient)
    txs:=result.Txs
    assert.Nil(t, err)
    assert.NotNil(t, txs)
    assert.True(t, len(txs) > 1)
    tx := txs[0]
    assert.True(t, len(tx.HashID) > 0)
    assert.Equal(t, 2, tx.BlockHeight)
}

func TestQueryAllBlocks(t *testing.T) {
    result, err := queryAllBlocks(0, dgClient)
    assert.Nil(t, err)
    blocks := result.Blocks
    assert.NotNil(t, blocks)
    //spew.Dump(blocks)
    assert.True(t, len(blocks) > 0, "blocks count should be more than 0")
}

func TestQueryByHash(t *testing.T) {
    result, _ := queryTxsInBlock(2, 0, 100, dgClient)
    txs := result.Txs
    tx := txs[0]
    result, err := queryByHash(tx.HashID[:3], dgClient)
    assert.Nil(t, err)
    assert.NotNil(t, result.Txs)
    assert.Equal(t, tx.HashID, result.Txs[0].HashID)
    //spew.Dump(result)

    result, _ = queryAllBlocks(1, dgClient)
    blocks:=result.Blocks
    block := blocks[0]
    result, err = queryByHash(block.HashID[:3], dgClient)
    assert.Nil(t, err)
    assert.NotNil(t, result.Blocks)
    assert.Equal(t, block.HashID, result.Blocks[0].HashID)
    //spew.Dump(result)
}

func TestQueryBlockRange(t *testing.T) {
    result, err := queryBlockRange(1, 3, dgClient)
    assert.Nil(t, err)
    assert.NotNil(t, result.Blocks)
    assert.Equal(t, 3, len(result.Blocks))
}

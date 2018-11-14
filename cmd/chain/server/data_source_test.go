package server

import (
    "encoding/json"
    "github.com/davecgh/go-spew/spew"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestSimulateDataSource_Block(t *testing.T) {
    dataSource := NewSimulateDataSource("ledger1")
    assert.Equal(t, 1, dataSource.LedgerHeight())
    dataSource.addNewBlock()

    assert.Equal(t, 2, dataSource.LedgerHeight())
    block := dataSource.Block(2)
    assert.NotNil(t, block)
    txs := dataSource.Transactions(2, 0, 2)
    assert.Equal(t, 2, len(txs))
    txs = dataSource.Transactions(2, 1, 2)
    assert.Equal(t, 1, len(txs))

    var record interface{}
    raw := txs[0].Content
    //zlog.PrettyJson(raw)
    err := json.Unmarshal(raw, &record)
    assert.Nil(t, err)
    spew.Dump(record)
}

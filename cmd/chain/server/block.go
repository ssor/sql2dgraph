package server

import (
    "fmt"
    "github.com/ssor/sql2graphql/helper"
)

type Block struct {
    Uid          string         `json:"uid,omitempty"`
    Hash         string         `json:"block_hash_id"`
    Height       int            `json:"block_height"`
    transactions []*Transaction `json:"-"`
    Ledger       *Ledger        `json:"in_ledger,omitempty"`
}

func newBlock(hashID string, height int, ledger *Ledger) *Block {
    return &Block{
        Hash:         hashID,
        Height:       height,
        Ledger:       ledger,
        transactions: []*Transaction{},
    }
}

func (block *Block) addTransactions(txs ...*Transaction) {
    for _, tx := range txs {
        tx.Block = block
    }
    block.transactions = append(block.transactions, txs...)
}

func (block *Block) DependentObjectHasUid() bool {
    if block.Ledger != nil {
        if len(block.Ledger.Uid) <= 0 {
            return false
        }
    }
    return true
}

func (block *Block) GetUidInfo() (string, string) {
    return fmt.Sprintf("block_%s", block.Hash), block.Uid
}

func (block *Block) QueryBy() (string, string) {
    return "block_hash_id", block.Hash
}

func (block *Block) SetUid(uid string) {
    block.Uid = uid
}

func (block *Block) Schemes() helper.Schemas {
    var schemes helper.Schemas
    schemes = schemes.Add(helper.NewSchemaIntIndex("block_height")).
        Add(helper.NewSchemaStringTrigramIndex("block_hash_id"))

    return schemes
}

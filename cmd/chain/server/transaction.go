package server

import (
    "crypto/sha256"
    "fmt"
    "github.com/ssor/sql2graphql/helper"
)

type Transaction struct {
    Uid          string `json:"uid,omitempty"`
    Hash         string `json:"tx_hash_id"`
    IndexInBlock int    `json:"index_in_block"`
    BlockHeight  int    `json:"tx.block_height"`
    Content      []byte `json:"-"`
    Block        *Block `json:"tx.block,omitempty"`
    //Block        *Block `json:"tx_in_block,omitempty"`
}

func (tx *Transaction) Schemes() helper.Schemas {
    var schemes helper.Schemas
    schemes = schemes.
        Add(helper.NewSchemaInt("index_in_block")).
        Add(helper.NewSchemaIntIndex("tx.block_height")).
        Add(helper.NewSchemaStringTrigramIndex("tx_hash_id"))

    return schemes
}

func newTransaction(content []byte, blockHeight, indexInBlock int) *Transaction {
    h := sha256.New()
    h.Write(content)
    v := h.Sum(nil)
    tx := &Transaction{
        Content:      content,
        Hash:         fmt.Sprintf("%X", v),
        IndexInBlock: indexInBlock,
        BlockHeight:  blockHeight,
    }
    return tx
}

func (tx *Transaction) DependentObjectHasUid() bool {
    if tx.Block != nil {
        if len(tx.Block.Uid) <= 0 {
            return false
        }
    }
    return true
}

func (tx *Transaction) GetUidInfo() (string, string) {
    return fmt.Sprintf("tx_%s", tx.Hash), tx.Uid
}

func (tx *Transaction) QueryBy() (string, string) {
    return "tx_hash_id", tx.Hash
}

func (tx *Transaction) SetUid(uid string) {
    tx.Uid = uid
}

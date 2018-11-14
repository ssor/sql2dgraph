package server

import (
    "fmt"
    "github.com/ssor/sql2graphql/helper"
)

type Ledger struct {
    Uid    string `json:"uid,omitempty"`
    Hash   string `json:"ledger_hash_id"`
    Height int    `json:"ledger_height"`
}

func newLedger(hashID string) *Ledger {
    return &Ledger{
        Hash: hashID,
    }
}

func (ledger *Ledger) DependentObjectHasUid() bool {
    return true
}

func (ledger *Ledger) GetUidInfo() (string, string) {
    return fmt.Sprintf("ledger_%s", ledger.Hash), ledger.Uid
}

func (ledger *Ledger) QueryBy() (string, string) {
    return "ledger_hash_id", ledger.Hash
}

func (ledger *Ledger) SetUid(uid string) {
    ledger.Uid = uid
}

func (ledger *Ledger) Schemes() helper.Schemas {
    var schemes helper.Schemas
    schemes = schemes.Add(helper.NewSchemaInt("ledger_height")).
        Add(helper.NewSchemaStringTrigramIndex("ledger_hash_id"))

    return schemes
}

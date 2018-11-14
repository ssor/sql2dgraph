package server

import (
    "fmt"
    "github.com/ssor/sql2graphql/cmd/chain/pipeline"
    "github.com/ssor/sql2graphql/helper"
)

func newLedgerUpdate(dataSource LedgerDataSource) *LedgerUpdateStep {
    return &LedgerUpdateStep{
        ledgerSource: dataSource,
    }
}

type LedgerUpdateStep struct {
    ledgerSource LedgerDataSource
    legerID      string
}

func (lu *LedgerUpdateStep) Exec(request *pipeline.Request) *pipeline.Result {
    ledger := lu.ledgerSource.Ledger()
    if ledger != nil {
        lu.legerID = ledger.Hash
        _, err := helper.MutationObj(ledger, dgClient)

        return &pipeline.Result{
            Error: err,
        }
    } else {
        return &pipeline.Result{}
    }
}

func (lu *LedgerUpdateStep) WorkDescription() (s string) {
    s = fmt.Sprintf("ledger [%s] update", lu.legerID)
    return
}

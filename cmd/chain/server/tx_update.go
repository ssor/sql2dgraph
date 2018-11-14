package server

import (
    "fmt"
    "github.com/ssor/sql2graphql/cmd/chain/pipeline"

    "github.com/ssor/sql2graphql/helper"
)

func newTxsUpdateStep(height int, dataSource LedgerDataSource) *TxsUpdateStep {
    return &TxsUpdateStep{
        height:       height,
        ledgerSource: dataSource,
    }
}

type TxsUpdateStep struct {
    height       int
    ledgerSource LedgerDataSource
}

func (tu *TxsUpdateStep) Exec(request *pipeline.Request) *pipeline.Result {
    totalCount := tu.ledgerSource.TransactionsCount(tu.height)
    if totalCount <= 0 {
        return &pipeline.Result{}
    }

    for i := 0; i < totalCount; i++ {
        txs := tu.ledgerSource.Transactions(tu.height, i, 1)
        if txs == nil || len(txs) <= 0 {
            return &pipeline.Result{
                Error: fmt.Errorf("block %d declare has %d txs, but in fact not, and tx of index %d not found", tu.height, totalCount, i),
            }
        }

        tx := txs[0]
        uid, err := helper.GetUidByIndex(tx.Block, dgClient)
        if err != nil {
            logger.Failedf("failed to get uid for block %d", tx.Block.Height)
            return &pipeline.Result{
                Error: err,
            }
        }

        tx.Block.SetUid(uid)
        _, err = helper.MutationObj(tx, dgClient)
        if err != nil {
            return &pipeline.Result{
                Error: fmt.Errorf("cannot mutate tx for %s", err),
            }
        }
    }
    return &pipeline.Result{}
}

func (tu *TxsUpdateStep) WorkDescription() (s string) {
    s = fmt.Sprintf("txs in block [%d] update", tu.height)
    return
}

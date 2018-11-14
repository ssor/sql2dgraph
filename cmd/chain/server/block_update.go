package server

import (
    "fmt"
    "github.com/ssor/sql2graphql/cmd/chain/pipeline"
    "github.com/ssor/sql2graphql/helper"
)

func newBlockUpdateStep(height int, dataSource LedgerDataSource) *BlockUpdateStep {
    return &BlockUpdateStep{
        height:       height,
        ledgerSource: dataSource,
    }
}

// add block to dgraph
type BlockUpdateStep struct {
    height       int
    ledgerSource LedgerDataSource
}

func (bu *BlockUpdateStep) Exec(request *pipeline.Request) *pipeline.Result {
    block := bu.ledgerSource.Block(bu.height)
    if block != nil {
        uid, err := helper.GetUidByIndex(block.Ledger, dgClient)
        if err != nil {
            logger.Failedf("failed to get uid for block %d", block.Height)
            return &pipeline.Result{
                Error: err,
            }
        }
        block.Ledger.SetUid(uid)
        _, err = helper.MutationObj(block, dgClient)

        return &pipeline.Result{
            Error: err,
        }
    } else {
        return &pipeline.Result{}
    }
}

func (bu *BlockUpdateStep) WorkDescription() (s string) {
    s = fmt.Sprintf("block [%d] update", bu.height)
    return
}

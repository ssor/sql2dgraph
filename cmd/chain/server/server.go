package server

import (
    "github.com/ssor/sql2graphql/cmd/chain/pipeline"
    "github.com/ssor/zlog"
)

var (
    serverLogger = zlog.New("chain", "stream", "server")
)

func CreateLegerUpdateServer(dataSource LedgerDataSource) *LedgerUpdateServer {
    server := &LedgerUpdateServer{
        chanAddBlock: make(chan int, 16),
        dataSource:   dataSource,
        workflow:     pipeline.New("wf_default"),
    }
    server.addWorkForAlterSchemas()
    server.addWorkForLedger()
    return server
}

type LedgerUpdateServer struct {
    chanAddBlock chan int
    dataSource   LedgerDataSource
    workflow     *pipeline.Pipeline
}

func (server *LedgerUpdateServer) AddBlockTask(height int) bool {
    if server.workflow.WorkCount() >= 3 {
        return false
    }
    server.addWorkForNewBlock(height)
    //server.chanAddBlock <- height
    return true
}

func (server *LedgerUpdateServer) Run() {
    zlog.Infof("running server ")
    server.workflow.Run()
}

func (server *LedgerUpdateServer) addWorkForAlterSchemas() {
    serverLogger.Infof("create stage for schema")
    server.workflow.NewWork(newSchemaUpdateStep())
}

func (server *LedgerUpdateServer) addWorkForLedger() {
    serverLogger.Infof("create stage for ledger")

    server.workflow.NewWork(newLedgerUpdate(server.dataSource))
}

func (server *LedgerUpdateServer) addWorkForNewBlock(height int) {
    serverLogger.Infof("add block stage [%d]", height)

    server.workflow.NewWork(newBlockUpdateStep(height, server.dataSource))
    server.workflow.NewWork(newTxsUpdateStep(height, server.dataSource))
}

package server

import (
    "container/list"
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "github.com/ssor/zlog"
    "strconv"
    "time"
)

var (
    simulateLogger = zlog.New("chain", "stream", "simulate")
)

func NewSimulateDataSource(ledgerHashID string) *SimulateDataSource {
    source := &SimulateDataSource{
        ledgerHashID: ledgerHashID,
        blocks:       list.New(),
    }
    source.blocks.PushBack(newBlock("genesis", 1, newLedger(source.ledgerHashID)))
    return source
}

// Simulate generate blocks as time goes by
type SimulateDataSource struct {
    ledgerHashID string
    blocks       *list.List
    eventHandler func(height int)
}

func (source *SimulateDataSource) SubscribeBlockChangedEvent(handler func(height int)) {
    source.eventHandler = handler
}

func (source *SimulateDataSource) Run() {
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for {
            <-ticker.C
            source.addNewBlock()
        }
    }()
}

func (source *SimulateDataSource) addNewBlock() {
    h := sha256.New()
    h.Write([]byte(strconv.FormatInt(int64(source.blocks.Len()), 10)))
    height := source.blocks.Len() + 1
    nb := newBlock(fmt.Sprintf("%X", h.Sum(nil)), height, newLedger(source.ledgerHashID))

    var txs []*Transaction
    raw, _ := json.Marshal(newAppleRecord())
    txs = append(txs, newTransaction(raw, height, len(txs)))
    raw, _ = json.Marshal(newBananaRecord())
    txs = append(txs, newTransaction(raw, height, len(txs)))
    nb.addTransactions(txs...)

    source.blocks.PushBack(nb)

    simulateLogger.Successf("block count increase to %d", source.blocks.Len())
    if source.eventHandler != nil {
        source.eventHandler(source.LedgerHeight())
    }
}

func (source *SimulateDataSource) Ledger() *Ledger {
    return &Ledger{
        Hash: source.ledgerHashID,
    }
}

func (source *SimulateDataSource) LedgerHeight() int {
    return source.blocks.Len()
}

func (source *SimulateDataSource) TransactionsCount(height int) int {
    block := source.Block(height)
    if block == nil {
        return 0
    }
    return len(block.transactions)
}

func (source *SimulateDataSource) Block(height int) *Block {
    var e *list.Element
    e = source.blocks.Front()
    if e == nil {
        return nil
    }
    for {
        v := e.Value.(*Block)
        if v.Height > height {
            return nil
        }
        if v.Height < height {
            e = e.Next()
            continue
        }
        return v
    }
    return nil
}

func (source *SimulateDataSource) Transactions(height, from, count int) []*Transaction {
    block := source.Block(height)
    if block == nil {
        return nil
    }
    if block.transactions == nil {
        return nil
    }
    totalTransactionCount := len(block.transactions)
    if from >= totalTransactionCount {
        return nil
    }
    endIndex := from
    if (from + count) >= totalTransactionCount {
        endIndex = totalTransactionCount
    } else {
        endIndex = from + count
    }
    return block.transactions[from:endIndex]
}

type TraceRecord struct {
    Id      int       `json:"id"`
    Name    string    `json:"name"`
    AtTime  time.Time `json:"at_time"`
    Others  string    `json:"others"`
    IsValid bool      `json:"is_valid"`
}

func newBananaRecord() *TraceRecord {
    tr := &TraceRecord{
        Id:     2,
        Name:   "banana",
        AtTime: time.Now(),
        Others: "That is banana",
    }
    return tr
}

func newAppleRecord() *TraceRecord {
    tr := &TraceRecord{
        Id:     1,
        Name:   "apple",
        AtTime: time.Now(),
        Others: "This is apple",
    }
    return tr
}

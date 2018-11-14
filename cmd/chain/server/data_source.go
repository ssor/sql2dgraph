package server

type LedgerDataSource interface {
    Ledger() *Ledger
    LedgerHeight() int
    Block(height int) *Block
    Transactions(height, from, count int) []*Transaction
    TransactionsCount(height int) int
    SubscribeBlockChangedEvent(handler func(height int))
}

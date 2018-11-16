package handler

import (
    "strconv"
    "strings"
)

const (
    noQueryLogic         = -1
    queryBlockRangeLogic = 0
    queryTxInBlockLogic  = 1
)

func newQueryLogic(blockFrom, blockTo, txFrom, txTo int) *queryLogic {
    logic := &queryLogic{
        blockFrom: blockFrom,
        blockTo:   blockTo,
        txFrom:    txFrom,
        txTo:      txTo,
    }
    logic.test()
    return logic
}

type queryLogic struct {
    blockFrom, blockTo, txFrom, txTo int
    result                           int
}

// test return what query type should be used
// -1 no type
// 0 query single block by height
// 1 query block in range
// 2 query txs in block
func (logic *queryLogic) test() {
    if logic.txFrom >= 0 {
        // I think you want to query tx
        if logic.txTo < 0 {
            logic.txTo = logic.txFrom
        }

        if logic.blockFrom <= 0 {
            logic.result = noQueryLogic
            return
        }
        logic.result = queryTxInBlockLogic
        return
    }

    if logic.blockFrom > 0 {
        if logic.blockTo >= logic.blockFrom {
            if (logic.blockTo - logic.blockFrom) >= 60 {
                logic.blockTo = logic.blockFrom + 59
            }
        } else {
            logic.blockTo = logic.blockFrom
        }
        logic.result = queryBlockRangeLogic
        return
    }

    logic.result = noQueryLogic
    return
}

func parseQuery(s string) (from, to int) {
    from, to = -1, -1
    ss := strings.Split(s, "-")
    if len(ss) <= 0 {
        return
    }

    if len(ss[0]) > 0 {
        i, err := strconv.Atoi(ss[0])
        if err != nil {
            logger.Warnf("parse [%s] failed for: %s", ss[0], err)
        } else {
            from = i
        }
    }

    if len(ss) > 1 {
        if len(ss[1]) > 0 {
            i, err := strconv.Atoi(ss[1])
            if err != nil {
                logger.Warnf("parse [%s] failed for: %s", ss[1], err)
            } else {
                to = i
            }
        }
    }
    return
}

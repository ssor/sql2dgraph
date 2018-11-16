package handler

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestParseQuery(t *testing.T) {
    from, to := parseQuery("")
    assert.Equal(t, 0, from)
    assert.Equal(t, 0, to)

    from, to = parseQuery("1")
    assert.Equal(t, 1, from)
    assert.Equal(t, 0, to)

    from, to = parseQuery("1-")
    assert.Equal(t, 1, from)
    assert.Equal(t, 0, to)

    from, to = parseQuery("1-2")
    assert.Equal(t, 1, from)
    assert.Equal(t, 2, to)
}

type queryInstance struct {
    block, tx    string
    expectResult int
}

func (ins *queryInstance) String() string {
    return fmt.Sprintf("block=%s tx=%s -> %d", ins.block, ins.tx, ins.expectResult)
}

func TestQueryLogic(t *testing.T) {
    list := []*queryInstance{
        {"", "", noQueryLogic},
        {"1", "", queryBlockRangeLogic},
        {"1-", "", queryBlockRangeLogic},
        {"1-1", "", queryBlockRangeLogic},
        {"1-1", "1", queryTxInBlockLogic},
        {"1-1", "1-", queryTxInBlockLogic},
        {"1-1", "-1", queryBlockRangeLogic},
        {"0", "", noQueryLogic},
        {"1-", "", queryBlockRangeLogic},
        {"1-0", "", queryBlockRangeLogic},
        {"1-0", "0", queryTxInBlockLogic},
        {"1-0", "0-", queryTxInBlockLogic},
        {"1-0", "-0", queryBlockRangeLogic},
    }
    for _, instance := range list {
        testQueryLogic(t, instance)
    }
}

func testQueryLogic(t *testing.T, instance *queryInstance) {
    blockFrom, blockTo := parseQuery(instance.block)
    txFrom, txTo := parseQuery(instance.tx)
    ql := newQueryLogic(blockFrom, blockTo, txFrom, txTo)
    assert.Equal(t, instance.expectResult, ql.result, instance.String())
}

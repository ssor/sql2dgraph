package handler

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func QueryAll(c *gin.Context) {
    keyword := c.Query("q")
    if len(keyword) >= 3 {
        logger.Info("query by keyword ...")
        status, res := queryByHashKeyword(keyword)
        c.JSON(status, res)
        return
    }

    block := c.Query("block")
    tx := c.Query("tx")
    blockFrom, blockTo := parseQuery(block)
    txFrom, txTo := parseQuery(tx)
    logic := newQueryLogic(blockFrom, blockTo, txFrom, txTo)
    status, res := queryByLogic(logic)
    c.JSON(status, res)
}

func queryByLogic(logic *queryLogic) (statusCode int, response *Response) {
    statusCode = http.StatusOK
    switch logic.result {
    //case 0:
    //    logger.Info("query single block ...")
    //    bi, err := querySingleBlockByHeight(logic.blockFrom, dgClient)
    //    if err != nil {
    //        logger.Failedf("query block by height [%d] failed for: %s", logic.blockFrom, err)
    //        response = newFailedResponse(noSearchResult)
    //        return
    //    }
    //    response = newSuccessResponse(bi)
    case queryBlockRangeLogic:
        logger.Info("query block range ...")
        result, err := queryBlockRange(logic.blockFrom, logic.blockTo, dgClient)
        if err != nil {
            logger.Failedf("query block in range [%d - %d] failed for: %s", logic.blockFrom, logic.blockTo, err)
            response = newFailedResponse(noSearchResult)
            return
        }
        response = newSuccessResponse(result)
    case queryTxInBlockLogic:
        logger.Info("query txs in block ...")
        result, err := queryTxsInBlock(logic.blockFrom, logic.txFrom, logic.txTo, dgClient)
        if err != nil {
            logger.Failedf("query tx in range [%d - %d] in block [%d] failed for: %s", logic.txFrom, logic.txTo, logic.blockFrom, err)
            response = newFailedResponse(noSearchResult)
            return
        }
        response = newSuccessResponse(result)
    default:
        logger.Failedf("no such result of %d", logic.result)
    }

    return
}

func queryByHashKeyword(kw string) (statusCode int, response *Response) {
    if len(kw) < 3 {
        statusCode = http.StatusBadRequest
        response = newFailedResponse(keywordTooShort)
        return
    }

    result, err := queryByHash(kw, dgClient)
    if err != nil {
        logger.Failedf("query hash failed for: %s", err)
        statusCode = http.StatusOK
        response = newFailedResponse(noSearchResult)
        return
    }
    statusCode = http.StatusOK
    response = newSuccessResponse(result)
    return
}

func QueryByHash(c *gin.Context) {
    id := c.Param("kw")
    if len(id) < 3 {
        c.JSON(http.StatusBadRequest, newFailedResponse(keywordTooShort))
        return
    }
    result, err := queryByHash(id, dgClient)
    if err != nil {
        logger.Failedf("query hash failed for: %s", err)
        c.JSON(http.StatusOK, newFailedResponse(noSearchResult))
        return
    }
    c.JSON(http.StatusOK, newSuccessResponse(result))
}

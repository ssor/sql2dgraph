package handler

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

const (
    keywordTooShort = "keyword must be more than three characters"
    noSearchResult  = "no search result"
)

type Response struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func newSuccessResponse(data interface{}) *Response {
    return &Response{
        Message: "OK",
        Data:    data,
    }
}

func newFailedResponse(msg string) *Response {
    return &Response{
        Message: msg,
    }
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

package handler

const (
    keywordTooShort = "keyword must be more than three characters"
    noSearchResult  = "no search result"
    paraError       = "参数错误"
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

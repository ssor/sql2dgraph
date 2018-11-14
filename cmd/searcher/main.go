package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ssor/sql2graphql/cmd/searcher/handler"
    "net/http"
)

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.GET("api/v1/search_hash/:kw", handler.QueryByHash)
    r.Static("/css", "./static/css")
    r.Static("/js", "./static/js")
    r.LoadHTMLGlob("templates/*")

    //views
    r.GET("/", index)
    r.GET("/index", index)
    r.GET("/search", search)

    r.Run("127.0.0.1:8081") // listen and serve on 0.0.0.0:8080
}
func index(context *gin.Context) {
    context.HTML(http.StatusOK, "index.tmpl", nil)
}
func search(context *gin.Context) {
    kw := context.Query("kw")
    context.HTML(http.StatusOK, "result.tmpl", gin.H{
        "keyword": kw,
    })
}

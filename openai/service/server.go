package service

import (
	"github.com/gin-gonic/gin"
	"log"
)

var accounts = gin.Accounts{
	"zxy": "d7yfhn85NArvatOzfz+i",
	"wq":  "TrHzU9dzwWTKlUvOUqAW0TVEMNqGrM",
}

const (
	listenerAddr = "0.0.0.0:9527"
)

func StartServer() {
	r := createRouter()
	err := r.Run(listenerAddr)
	if err != nil {
		log.Panic(err)
	}
}

func createRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1", gin.BasicAuth(accounts))
	{
		v1.POST("/ping", ping)
		v1.POST("/chat", chatHandler)
		v1.POST("/image", imageHandler)
	}
	return router
}

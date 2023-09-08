package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/registry/verify", verify)
	router.GET("/v2/registry/verify/manifests/latest", test)
	router.GET("/v2", v2)
	router.POST("/v2/users/login", login)
	router.Run(":8888")
}

func login(c *gin.Context) {
	fmt.Println(c.Params)
}

func v2(c *gin.Context) {
	c.Header("WWW-Authenticate", `Bearer realm="http://zxy.kingqian.wang:42010/token",service="zxy.kingqian.wang"`)
	c.AbortWithStatus(401)
}

func verify(c *gin.Context) {
	fmt.Println(c.Request.URL)
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.Header("WWW-Authenticate", `Basic realm="Restricted Access!"`)
		c.AbortWithStatus(401)
		return
	}
	c.JSON(200, map[string]string{})
}

func test(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	fmt.Println("Authorization header:", authHeader)
	fmt.Println(c.Request.URL)
}

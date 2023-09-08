package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		fmt.Println("start_hello")

		commonName := r.Header.Get("X-SSL-CLIENT-S-DN")
		sn := r.Header.Get("X-SSL-CLIENT-SERIAL")
		writer.Write([]byte("HelloWorld\ncommon_name:" + commonName + "\n serial number:" + sn))
	})
	http.ListenAndServe(":1234", nil)

	r := gin.Default()
	r.GET("/", func(context *gin.Context) {

	})

}

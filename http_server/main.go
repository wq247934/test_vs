package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", hello)
	if err := http.ListenAndServe("127.0.0.1:12345", nil); err != nil {
		fmt.Println(err)
	}

}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

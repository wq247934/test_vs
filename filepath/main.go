package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	path := "/root/work/release/www.baidu.com/dns01"
	fmt.Println(filepath.Base(path))
	fmt.Println(filepath.Base(filepath.Dir(path)))
}

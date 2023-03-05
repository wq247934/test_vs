package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	//var files []string
	//
	//root := ""
	//err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
	//	stat, _ := os.Stat(path)
	//	if stat.IsDir() {
	//		return nil
	//	}
	//	files = append(files, path)
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}
	//for _, file := range files {
	//	fmt.Println(file)
	//}

	//fmt.Println(filepath.Join("/etc", "nginx/", "conf.d"))
	path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("path:", path)
	fmt.Println("dir path:", filepath.Dir(path))
}

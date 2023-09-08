package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	req, err := http.NewRequest("GET", "https://item.jd.com/10039842742321.html", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`img src="(.+?)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		fmt.Println(match[1])
		download(fmt.Sprintf("https:%s", match[1]))
	}
}

func download(path string) {
	resp, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filepath.Base(path), body, 0600)
	if err != nil {
		panic(err)
	}
}

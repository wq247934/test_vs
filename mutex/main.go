package main

import (
	"log"
	"sync"
	"time"
)

var m = make(map[int]int)

func main() {
	log.Println()
	m[1] = 1
	m[2] = 2
	m[3] = 3
	go get(1)
	mu := sync.RWMutex{}
	mu.Lock()

	modify()
	time.Sleep(time.Second * 3)
	mu.Unlock()
	time.Sleep(time.Second * 5)
	log.Println()
}

func modify() {
	m[1] = 111111
}

func get(key int) int {
	time.Sleep(time.Second * 1)
	log.Println(m[key])
	return m[key]
}

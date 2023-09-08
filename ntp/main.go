package main

import (
	"fmt"
	"github.com/beevik/ntp"
)

func main() {
	time, err := ntp.Time("0.debian.pool.ntp.org")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.String())
}

package main

import (
	"fmt"
	"github.com/robfig/cron"
	"time"
)

func main() {
	c := cron.New()
	err := c.AddFunc("00 * * * * * ", func() {
		fmt.Println(time.Now())
	})
	if err != nil {
		fmt.Println(err)
	}

	c.Start()
	select {}
}

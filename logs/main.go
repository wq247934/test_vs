package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type connectLog struct {
	User       string
	RemoteAddr string
	Connected  bool
}

func (log connectLog) String() string {
	if log.Connected {
		return fmt.Sprintf("%s connected to the server,orion ip is %s", log.User, log.RemoteAddr)
	} else {
		return fmt.Sprintf("%s disconnected to the server,orion ip is %s", log.User, log.RemoteAddr)
	}

}

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	file, err := os.OpenFile("./a.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println(err)
	}
	log1 := connectLog{
		User:       "wangqian",
		RemoteAddr: "192.168.1.1",
		Connected:  true,
	}
	log2 := connectLog{
		User:       "wangqian",
		RemoteAddr: "192.168.1.1",
		Connected:  false,
	}
	logger.SetOutput(file)
	logger.Info(log1)
	logger.Info(log2)
}

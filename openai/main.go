package main

import (
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"test_vs/openai/chat"
	"test_vs/openai/image"
	"test_vs/openai/service"
)

var (
	sk     string
	url    string
	client *openai.Client
)

func init() {
	v := viper.GetViper()
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	sk = v.GetString("openai.sk")
	url = v.GetString("openai.url")
	config := openai.DefaultConfig(sk)
	config.BaseURL = url
	client = openai.NewClientWithConfig(config)
}

func main() {
	chat.InitChatClient(client)
	image.InitImageClient(client)
	chat.StartClearHistoryJob()
	service.StartServer()
}

package main

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

func main() {
	v := viper.GetViper()
	v.SetConfigFile("config.yaml")
	v.GetString("openai.sk")
	config := openai.DefaultConfig(v.GetString("openai.sk"))
	config.BaseURL = v.GetString("openai.base_url")
	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

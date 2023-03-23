package chat_gpt

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"time"
)

var (
	sk              string
	url             string
	client          *openai.Client
	userChatHistory = make(map[string]*chatHistory)
	systemContent   = "你是一个语言模型AI助手，专门为您提供帮助和回答问题的。你可以回答各种问题，包括但不限于科学、历史、文化、生活、体育、技术和其他相关领域的问题"
)

type chatHistory struct {
	username     string
	history      []openai.ChatCompletionMessage
	lastChatTime time.Time
}

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
	clearHistory()
}

func Chat(username, content string) string {
	history, ok := userChatHistory[username]
	if !ok {
		history = &chatHistory{
			username:     username,
			history:      []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: systemContent}},
			lastChatTime: time.Now(),
		}
		userChatHistory[username] = history
	}

	messages := append(history.history, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: content})
	var result string
	errCount := 0
	for {
		var err error
		result, err = CreateChatCompletion(messages)
		if err != nil {
			fmt.Println(err)
			if errCount > 3 {
				return err.Error()
			}
			errCount++
			continue
		}
		break
	}
	messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: result})
	history.history = messages
	history.lastChatTime = time.Now()
	return result
}

func CreateChatCompletion(messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return "", errors.New(fmt.Sprintf("ChatCompletion error: %v\n", err))
	}

	return resp.Choices[0].Message.Content, nil
}

func clearHistory() {
	c := cron.New()
	//每30分钟清理一次上下文
	c.AddFunc("*/30 * * * *", func() {
		for _, history := range userChatHistory {
			if time.Now().Sub(history.lastChatTime).Minutes() > 30 {
				history.history = []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: systemContent}}
			}
		}
	})
	c.Start()
}

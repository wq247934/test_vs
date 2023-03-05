package main

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func main() {
	c := gogpt.NewClient("")
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3Dot5Turbo,
		MaxTokens: 5,
		Prompt:    "Lorem ipsum",
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Choices[0].Text)
}

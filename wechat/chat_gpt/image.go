package chat_gpt

import (
	"context"
	"encoding/base64"
	"github.com/sashabaranov/go-openai"
)

func CreateImage(prompt string) ([]byte, error) {
	imageResp, err := client.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Prompt:         prompt,
			N:              1,
			Size:           "1024x1024",
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		})
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(imageResp.Data[0].B64JSON)
}

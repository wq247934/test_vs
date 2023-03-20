package chat_gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

var sk, url string

func init() {
	config := viper.New()
	config.SetConfigFile("config.yaml")
	sk = config.GetString("openai.sk")
	url = config.GetString("openai.url")
	if sk == "" || url == "" {
		panic("please check config")
	}
}

var userMap = make(map[string]string)

func Chat(username, content string) string {
	sessionID, ok := userMap[username]
	if !ok {
		sessionID = uuid.New().String()
		userMap[username] = sessionID
	}
	fmt.Printf("username:%s,content:%s,sessionID:%s\n", username, content, sessionID)
	data, err := postGpt(sessionID, content)
	if err != nil {
		fmt.Println(err)
		return "出现错误,请联系管理员"
	}
	return data
}

func Reset(username string) string {
	delete(userMap, username)
	return "对话已重置"
}

type Payload struct {
	APIKey    string `json:"apiKey"`
	SessionID string `json:"sessionId"`
	Content   string `json:"content"`
}

func postGpt(sessionID, content string) (string, error) {
	data := Payload{
		APIKey:    sk,
		SessionID: sessionID,
		Content:   content,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		// handle error
	}

	payload := bytes.NewReader(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	resp := make(map[string]string)
	json.Unmarshal(body, &resp)
	return resp["data"], nil
}

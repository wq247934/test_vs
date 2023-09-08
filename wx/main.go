package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	Token       = "wechat_token"
	BingChatAPI = "http://ssh.kingqian.wang:5001"
)

func main() {
	r := gin.Default()
	r.Use(verify())
	r.GET("/wx", func(c *gin.Context) {
		echostr := c.Query("echostr")
		_, err := c.Writer.Write([]byte(echostr))
		if err != nil {
			panic(err)
		}
	})
	r.GET("/hello", func(c *gin.Context) {
		c.Writer.Write([]byte("world"))
	})
	r.POST("/wx", WXMsgReceive)
	r.Run(":8081")
}

type WXRepTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}

// WXTextMsg 微信文本消息结构体
type WXTextMsg struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgId        int64  `xml:"MsgId"`
}

func WXMsgReceive(c *gin.Context) {
	var textMsg WXTextMsg
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}

	log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, textMsg.Content)
	WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, textMsg.CreateTime)
}

func WXMsgReply(c *gin.Context, fromUser, toUser string, t int64) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   t,
		MsgType:      "text",
		Content:      fmt.Sprintf("[消息回复] - %s", time.Now().Format("2006-01-02 15:04:05")),
	}
	time.Sleep(time.Second * 5)
	c.XML(200, repTextMsg)
}

func verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		list := []string{Token, timestamp, nonce}
		sort.Strings(list)
		hash := sha1.New()
		hash.Write([]byte(strings.Join(list, "")))
		hashcode := hex.EncodeToString(hash.Sum(nil))
		if hashcode != signature {
			c.AbortWithError(http.StatusBadRequest, errors.New("InvalidSignature"))
			logrus.Errorln("Failed to verify signature")
			return
		}
		logrus.Info("verify signature success!")
		c.Next()
	}
}

type bingChat struct {
	Style    string
	Question string
	Token    string
}

type bingReply struct {
	answer   string
	suggests []string
	urls     []string
	reset    bool
	token    string
}
type Resp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (c *bingChat) chat() (*bingReply, error) {
	var r Resp
	client := req.C()
	_, err := client.R().
		SetSuccessResult(&r).
		Get(fmt.Sprintf("%s/style=%s&question=%s", BingChatAPI, c.Style, c.Question))
	if err != nil {
		return nil, err
	}
	client.OnAfterResponse(handleResponse)
	data := r.Data.(bingReply)
	return &data, nil
}

func handleResponse(c *req.Client, r *req.Response) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var resp Resp

	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}
	if resp.Code != 200 {
		return errors.New(resp.Message)
	}
	return nil
}

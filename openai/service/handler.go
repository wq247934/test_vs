package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"test_vs/openai/chat"
)

func imageHandler(c *gin.Context) {

}

func chatHandler(c *gin.Context) {
	body := make(map[string]interface{})
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"code":    "ParseJsonError",
			"message": err.Error(),
		})
		return
	}
	user, ok := body["user"].(string)
	if !ok || strings.TrimSpace(user) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"code":    "userIsEmpty",
			"message": "user cannot be empty.",
		})
		return
	}
	prompt := body["prompt"].(string)
	if !ok || strings.TrimSpace(prompt) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"code":    "PromptIsEmpty",
			"message": "prompt cannot be empty.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"code":   "code",
		"data":   chat.Chat(user, prompt),
	})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"code":   "ok",
		"data":   "pong",
	})
}

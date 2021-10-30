package routers

import (
	news_bot "github.com/explabs/ad-ctf-paas-api/news-bot"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func NewsHandler(c *gin.Context) {
	tb := news_bot.TelegramBot{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:           os.Getenv("CHAT_ID"),
		NewsFolder:       "",
		Round:            0,
	}
	if err := tb.LoadMessage("test.md"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err,
		})
		return
	}
	if err := tb.SendMessage(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

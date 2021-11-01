package routers

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/database"
	news_bot "github.com/explabs/ad-ctf-paas-api/news-bot"
	"github.com/explabs/ad-ctf-paas-api/walker/providers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func NewsHandler(c *gin.Context) {
	var news providers.RoundsStruct
	news.Parse("exploits.yml")
	round, _ := database.GetRound()
	if round < len(news.Rounds) {
		if err := news_bot.SendNews([]string{
			news.Rounds[round].HintNews,
			news.Rounds[round].News,
		}); err != nil{
			log.Println(err)
		}
		c.Data(http.StatusOK, "plain/text", []byte(fmt.Sprintf("news{round=%d} 1", round)))
	}
	c.Data(http.StatusOK, "plain/text", []byte(fmt.Sprintf("news{round=%d} 0", round)))
}

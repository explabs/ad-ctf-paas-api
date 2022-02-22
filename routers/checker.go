package routers

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/rabbit"
	"github.com/explabs/ad-ctf-paas-api/walker"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CheckerHandler(c *gin.Context) {
	checkResult, err := walker.CheckFlags()
	if err != nil {
		log.Println(err)
		c.Data(http.StatusOK, "text/plain", []byte(""))
	}
	var data string
	for k, v := range checkResult {
		data += fmt.Sprintf("%s %d\n", k, v)
	}
	database.RemoveAllFlags()
	database.WriteTime()
	putResult, err := walker.PutFlags()
	if err != nil {
		log.Println(err)
		c.Data(http.StatusOK, "text/plain", []byte(""))
	}
	for k, v := range putResult {
		data += fmt.Sprintf("%s %d\n", k, v)
	}

	c.Data(http.StatusOK, "text/plain", []byte(data))
}

func RunChecker(c *gin.Context) {
	if err := rabbit.SendMessage("checker", "{\"type\": \"start\"}"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"checker": "started"})
}

func StopChecker(c *gin.Context) {
	if err := rabbit.SendMessage("checker", "{\"type\": \"stop\"}"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"checker": "stopped"})
}

package routers

import (
	"github.com/explabs/ad-ctf-paas-api/config"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InfoHandler(c *gin.Context){
	regStatus, err := database.RegistrationStatus()
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"detail": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"mode":config.Conf.Mode,
		"reg_status": regStatus,
	})
}

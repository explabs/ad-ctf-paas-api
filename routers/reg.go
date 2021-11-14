package routers

import (
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/gin-gonic/gin"
	"net/http"
)
func OpenRegistrationHandler(c *gin.Context){
	err := database.ChangeRegistrationStatus("open")
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"detail": err})
	}
	InfoHandler(c)
}
func CloseRegistrationHandler(c *gin.Context){
	err := database.ChangeRegistrationStatus("close")
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"detail": err})
	}
	InfoHandler(c)
}

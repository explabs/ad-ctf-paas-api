package routers

import (
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"os"
)

func AddVpnTeam(team *models.TeamInfo) error {
	vpnAddr := os.Getenv("OVPN_ADMIN")
	if vpnAddr == ""{
		vpnAddr = "http://localhost:9000"
	}
	urlAddr := vpnAddr + "/api/user/create"
	_, httpErr := http.PostForm(urlAddr, url.Values{
		"username": {team.Name},
		"password": {"kb4ctf"},
	})

	if httpErr != nil{
		log.Println(httpErr)
		return httpErr
	}
	return nil
}

func CreateVpnTeams(c *gin.Context){
	teams, dbErr := database.GetTeams()
	if dbErr != nil{
		c.JSON(http.StatusBadRequest, dbErr)
		return
	}
	for _, team := range teams{
		AddVpnTeam(team)
	}
}
package routers

import (
	"bytes"
	"encoding/json"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

var vpnUrl = "http://openvpn:9000/"

func AddVpnTeam(team *models.Team, rawPassword string) error {
	urlAddr := vpnUrl + "api/user/create"
	_, httpErr := http.PostForm(urlAddr, url.Values{
		"username": {team.Login},
		"password": {rawPassword},
	})
	if httpErr != nil {
		log.Println(httpErr)
		return httpErr
	}
	return nil
}

type VpnRoute struct {
	User          string         `json:"User"`
	ClientAddress string         `json:"ClientAddress"`
	CustomRoutes  []CustomRoutes `json:"CustomRoutes"`
}
type CustomRoutes struct {
	Address string `json:"Address"`
	Mask    string `json:"Mask"`
}

func DownloadVpnConfig(username string) (string, error) {
	urlAddr := vpnUrl + "api/user/config/show"
	response, httpErr := http.PostForm(urlAddr, url.Values{
		"username": {username},
	})
	if httpErr != nil {
		log.Println(httpErr)
		return "", httpErr
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	vpnConfig := string(responseData)

	return vpnConfig, nil
}

func (vpnRoute *VpnRoute) WriteTeamsRoutes() error {
	urlAddr := vpnUrl + "api/user/ccd/apply"

	jsonValue, _ := json.Marshal(vpnRoute)
	_, httpErr := http.Post(urlAddr, "application/json", bytes.NewBuffer(jsonValue))
	if httpErr != nil {
		log.Println(httpErr)
		return httpErr
	}
	return nil
}

func GetVpnConfigHandler(c *gin.Context) {
	user, _ := c.Get("id")
	username := user.(*models.JWTTeam).TeamName
	log.Println(user, username)
	vpnConfig, err := DownloadVpnConfig(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err})
		return
	}
	c.Data(200, "plain/text; charset=utf-8", []byte(vpnConfig))

}

func AddVpnRoutes(c *gin.Context) {
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": dbErr})
		return
	}

	for _, team := range teams {
		_, ipNet, _ := net.ParseCIDR(team.Address)
		route := CustomRoutes{
			Address: ipNet.IP.String(),
			Mask:    net.IP(ipNet.Mask).String(),
		}
		vpnRoute := VpnRoute{
			User:          team.Login,
			ClientAddress: "dynamic",
			CustomRoutes:  []CustomRoutes{route},
		}
		if err := vpnRoute.WriteTeamsRoutes(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": err})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vpn routes added for users"})

}

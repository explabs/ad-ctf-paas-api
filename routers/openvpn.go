package routers

import (
	"bytes"
	"encoding/json"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

func AddVpnTeam(team *models.Team, rawPassword string) error {
	vpnAddr := os.Getenv("OVPN_ADMIN")
	if vpnAddr == "" {
		vpnAddr = "http://localhost:9000"
	}
	urlAddr := vpnAddr + "/api/user/create"
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
	User          string       `json:"User"`
	ClientAddress string       `json:"ClientAddress"`
	CustomRoutes  []CustomRoutes `json:"CustomRoutes"`
}
type CustomRoutes struct {
	Address string `json:"Address"`
	Mask    string `json:"Mask"`
}

func (vpnRoute *VpnRoute) WriteTeamsRoutes() error {
	vpnAddr := os.Getenv("OVPN_ADMIN")
	if vpnAddr == "" {
		vpnAddr = "http://localhost:9000"
	}
	urlAddr := vpnAddr + "/api/user/ccd/apply"

	jsonValue, _ := json.Marshal(vpnRoute)
	_, httpErr := http.Post(urlAddr, "application/json", bytes.NewBuffer(jsonValue))
	if httpErr != nil {
		log.Println(httpErr)
		return httpErr
	}
	return nil
}

func CreateVpnTeams(c *gin.Context) {
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, dbErr)
		return
	}

	for _, team := range teams {
		_, ipNet, _ := net.ParseCIDR(team.Address)
		route := CustomRoutes{
			Address: ipNet.IP.String(),
			Mask: net.IP(ipNet.Mask).String(),
		}
		vpnRoute := VpnRoute{
			User:          team.Login,
			ClientAddress: "dynamic",
			CustomRoutes: []CustomRoutes{route},
		}
		vpnRoute.WriteTeamsRoutes()
	}
}

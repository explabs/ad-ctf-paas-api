package routers

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func TeamsList(c *gin.Context) {
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"teams": teams,
	})
}

func DeleteTeams(c *gin.Context) {
	team := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("team %s deleted", team),
	})
}

type TeamsForAnsible struct {
	IP        string `json:"ip"`
	Netmask      string `json:"netmask"`
	Mode      string `json:"mode"`
	Name      string `json:"name"`
	DHCPStart string `json:"dhcp_start"`
	DHCPEnd   string `json:"dhcp_end"`
}

// TODO: fix this func if it will be used
//func CountTeamsHandler(c *gin.Context) {
//	var result []TeamsForAnsible
//	teams, _ := database.GetTeams()
//	for _, team := range teams{
//		data := TeamsForAnsible{
//			Name: team.Name,
//			Netmask: "255.255.255.0",
//			Mode: "nat",
//		}
//		// network address
//		ip := net.ParseIP(team.Address)
//		ip = ip.To4()
//		ip[3] = 0
//		data.IP = ip.String()
//		// dhcp start
//		ip[3] = 11
//		data.DHCPStart = ip.String()
//		ip[3] = 253
//		data.DHCPEnd = ip.String()
//		result = append(result, data)
//	}
//
//	c.JSON(http.StatusOK, gin.H{"teams": result})
//}

func prometheusManagerRequest(action string) (string, error) {
	url := "http://localhost:9091/"
	switch action {
	case "start":
		url = url + action
	case "stop":
		url = url + action
	default:
		return "",fmt.Errorf("bad action: %s", action)
	}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("admin", os.Getenv("ADMIN_PASS"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	responseText, err := ioutil.ReadAll(resp.Body)
	return string(responseText), err
}


func RunPrometheusHandler(c *gin.Context){
	response, err := prometheusManagerRequest("start")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": response})
}

func StopPrometheusHandler(c *gin.Context){
	response, err := prometheusManagerRequest("stop")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": response})
}

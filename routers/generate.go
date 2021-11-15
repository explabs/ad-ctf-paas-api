package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/explabs/ad-ctf-paas-api/config"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/pkg/archive"
	"github.com/explabs/ad-ctf-paas-api/pkg/temporary"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Teams struct {
	Name  string
	OSImg string
}

func SshKeyArchiveHandler(c *gin.Context) {
	archiveName := "keys.tar.gz"
	defer os.RemoveAll(archiveName)
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}

	tmpdirName, _ := temporary.CreateTempDir()
	defer os.RemoveAll(tmpdirName)
	for _, team := range teams {
		fileName := fmt.Sprintf("%s.pub", team.Name)
		temporary.WriteFileDataToDir(tmpdirName, fileName, team.SshPubKey)
	}
	archive.Compress(tmpdirName, archiveName)
	c.File(archiveName)
}

func GenerateVariables(c *gin.Context) {
	fileName := "teams.tf"

	// osImageFilename := "focal-server-cloudimg-amd64.img"
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}
	hclFile := hclwrite.NewEmptyFile()

	// create new file on system
	tfFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tfFile.Close()
	defer os.Remove(fileName)
	// initialize the body of the new file object
	rootBody := hclFile.Body()

	// generate os_image variable
	vmBlock := rootBody.AppendNewBlock("variable", []string{"teams"})
	vmBlockBody := vmBlock.Body()
	var teamsList []cty.Value
	var ipsList []cty.Value
	var cidrList []cty.Value
	for _, team := range teams {
		teamsList = append(teamsList, cty.StringVal(team.Login))
		teamIP, teamNetIp, _ := net.ParseCIDR(team.Address)
		ipsList = append(ipsList, cty.StringVal(teamIP.To4().String()))
		cidrList = append(cidrList, cty.StringVal(teamNetIp.String()))
	}
	vmBlockBody.SetAttributeValue("default", cty.ListVal(teamsList))

	ips := rootBody.AppendNewBlock("variable", []string{"ips"})
	ipsBody := ips.Body()
	ipsBody.SetAttributeValue("default", cty.ListVal(ipsList))

	cidrs := rootBody.AppendNewBlock("variable", []string{"cidr"})
	cidrsBody := cidrs.Body()
	cidrsBody.SetAttributeValue("default", cty.ListVal(cidrList))
	hclFile.WriteTo(tfFile)
	c.File(fileName)
}

// curl -u admin:admin -H "Content-Type: application/json" --data '{ "password": "test", "target": "192.168.100.105:8080", "interval": "30s", "timeout": "10s", "jobs": [{"name": "checker", "path": "api/v1/game/checker"}]}' http://localhost:9091/generate
// TOKEN=$(curl -s -X POST -H 'Accept: application/json' -H 'Content-Type: application/json' --data '{"username":"admin","password":"admin"}' http://localhost/api/v1/auth/login | jq -r '.token')
// curl -H 'Accept: application/json' -H "Authorization: Bearer ${TOKEN}" -X POST http://localhost/api/v1/admin/generate/prometheus
type Payload struct {
	Password string `json:"password"`
	Target   string `json:"target"`
	Interval string `json:"interval"`
	Timeout  string `json:"timeout"`
	Jobs     []Jobs `json:"jobs"`
}
type Jobs struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Interval string `json:"interval"`
}

func GeneratePrometheus(c *gin.Context) {
	interval := config.Conf.RoundInterval
	data := Payload{
		Password: config.Conf.CheckerPassword,
		Target:   "ad-api",
		Interval: interval,
		Timeout:  "10s",
	}

	baseApiPath := "api/v1/game/"
	defenceJobNames := []string{"checker", "exploits", "news"}
	for _, name := range defenceJobNames {
		data.Jobs = append(data.Jobs, Jobs{Name: name, Path: baseApiPath + name, Interval: interval})
		if config.Conf.Mode != "defence" {
			break
		}
		interval = config.Conf.ExploitInterval
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://prometheus-manager:9091/generate", body)
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth("admin", os.Getenv("ADMIN_PASS"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}

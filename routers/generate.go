package routers

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/pkg/archive"
	"github.com/explabs/ad-ctf-paas-api/pkg/temporary"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"net"
	"net/http"
	"os"
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

func GeneratePrometheus(c *gin.Context) {

}

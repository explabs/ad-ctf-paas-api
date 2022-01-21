package routers

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/config"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"time"
)

type Team struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	SshPubKey string `json:"ssh_pub_key"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetTeamInfo(c *gin.Context) {
	team, err := database.GetTeam("")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"team": team[0]})
}

func generateIp(number int) string {
	ip, _, err := net.ParseCIDR(config.Conf.Network)
	if err != nil {
		log.Println(err.Error())
	}
	ip = ip.To4()
	ip[2] = byte(number)
	ip[3] = 10
	// TODO: find better solution for generate cidr
	return ip.String() + "/24"
}

func CreateTeam(c *gin.Context) {
	status, err := database.RegistrationStatus()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err})
		return
	}
	if status == "close" {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "registration closed"})
		return
	}
	var team Team
	jsonErr := c.BindJSON(&team)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": jsonErr.Error()})
		return
	}
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		log.Println(dbErr)
	}

	// check if user already exists
	for _, dbTeam := range teams {
		if dbTeam.Login == slug.Make(team.Name) {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "team already exists"})
			return
		}
	}

	ipAddress := generateIp(len(teams) + 1)
	hash, hashErr := HashPassword(team.Password)
	if hashErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": hashErr.Error()})
		return
	}

	dbTeam := &models.Team{
		ID:        primitive.NewObjectID(),
		Name:      team.Name,
		Login:     slug.Make(team.Name),
		Address:   ipAddress,
		Hash:      hash,
		SshPubKey: team.SshPubKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	dbErr = database.CreateTeam(dbTeam)
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": dbErr.Error()})
		return
	}

	if err := AddVpnTeam(dbTeam, team.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": hashErr.Error()})
		return
	}
	c.SetCookie("login",
		slug.Make(team.Name),
		60*60*24,
		"/",
		"localhost",
		false,
		false)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("The team %s created", team.Name),
	})
}

func DeleteTeam(c *gin.Context) {
	teamName := c.Param("name")
	dbErr := database.DeleteTeam(teamName)
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": dbErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s deleted", teamName),
	})
}

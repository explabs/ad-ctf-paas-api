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
	"strings"
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
	user, _ := c.Get("id")
	team, err := database.GetTeam(user.(*models.JWTTeam).TeamName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"team": team})
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

type CredPolicy struct {
	NameMaxLen     int
	NameMinLen     int
	PasswordMinLen int
	PasswordMaxLen int
}

func (c *CredPolicy) SetPolicy() {
	c.PasswordMinLen, c.PasswordMaxLen = 14, 40
	c.NameMinLen, c.NameMaxLen = 3, 40
}

func (c *CredPolicy) CheckCredPolicy(team *Team) error {
	if c.NameMinLen > len(team.Name) {
		return fmt.Errorf("short name: %s", team.Name)
	}

	return nil
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
	login := slug.Make(team.Name)
	login = strings.Replace(login, "-", "_", -1)
	dbTeam := &models.Team{
		ID:        primitive.NewObjectID(),
		Name:      team.Name,
		Login:     login,
		Address:   ipAddress,
		Hash:      hash,
		SshPubKey: team.SshPubKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	dbErr = database.CreateTeam(dbTeam)
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": dbErr})
		return
	}

	vpnErr := AddVpnTeam(dbTeam, team.Password)
	if vpnErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": vpnErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("The team %s created", team.Name),
		"login":   login,
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

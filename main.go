package main

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/explabs/ad-ctf-paas-api/config"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/middlewares"
	"github.com/explabs/ad-ctf-paas-api/models"
	"github.com/explabs/ad-ctf-paas-api/routers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func addAdmin() {
	teams, _ := database.GetUsers()
	for _, team := range teams {
		if team.Login == "admin" {
			return
		}
	}
	password := os.Getenv("ADMIN_PASS")
	if password == "" {
		password = "admin"
	}
	hash, _ := routers.HashPassword(password)
	database.CreateTeam(&models.Team{
		ID:        primitive.NewObjectID(),
		Login:     "admin",
		Hash:      hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())
	err := config.ReadConf("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	// init db connections
	database.InitMongo()
	database.InitRedis()

	// add info about services to mongo
	var sc config.ServicesInfo
	sc.Load()
	database.UploadServices(sc.Services)

	// add admin user to database
	addAdmin()

	// create router object
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// the jwt middleware
	authMiddleware, err := jwt.New(&middlewares.JwtMiddlewareStruct)

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(http.StatusNotFound, gin.H{"detail": "Not found"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.GET("/auth/refresh_token", authMiddleware.RefreshHandler)

		v1.Use(authMiddleware.MiddlewareFunc())
		v1.GET("/scoreboard", routers.ShowScoreboard)
		v1.GET("/scoreboard/:name", routers.ShowTeamStatus)
		if config.Conf.Mode == "attack-defence"{
			v1.POST("/submit", routers.SubmitFlagHandler)
		}

		// teams CRUD
		team := v1.Group("/team")
		team.POST("/", routers.CreateTeam)
		team.Use(authMiddleware.MiddlewareFunc())
		{
			team.GET("/", routers.GetTeamInfo)
			team.DELETE("/", routers.DeleteTeam)
		}

		// admins functions
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.MiddlewareFunc())
		admin.Use(middlewares.IsAdmin())
		{
			admin.GET("/teams", routers.TeamsList)
			admin.POST("/vpn", routers.CreateVpnTeams)
			admin.DELETE("/team/:name", routers.DeleteTeams)
			admin.POST("/generate/variables", routers.GenerateVariables)
			admin.POST("/generate/sshkeys", routers.SshKeyArchiveHandler)
			admin.POST("/generate/prometheus", routers.GeneratePrometheus)
			admin.POST("/prom/start", routers.RunPrometheusHandler)
			admin.POST("/prom/stop", routers.StopPrometheusHandler)
			admin.GET("/reg/open", routers.OpenRegistrationHandler)
			admin.GET("/reg/close", routers.CloseRegistrationHandler)
		}

		// public services without auth
		services := v1.Group("/services")
		{
			services.GET("/teams/info", routers.CountTeamsHandler)
			services.GET("/system/info", routers.InfoHandler)
		}

		// routes for prometheus with basic auth
		walker := v1.Group("/game")
		walker.Use(gin.BasicAuth(gin.Accounts{
			"checker": config.Conf.CheckerPassword,
		}))
		{
			walker.GET("/checker", routers.CheckerHandler)
			walker.GET("/news", routers.NewsHandler)
			walker.GET("/exploit", routers.ExploitHandler)
		}
	}
	router.Run(":8080")
}

package main

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/explabs/ad-ctf-paas-api/config"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/middlewares"
	"github.com/explabs/ad-ctf-paas-api/models"
	"github.com/explabs/ad-ctf-paas-api/routers"
	"github.com/gin-contrib/cors"
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
	// load config data
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type","Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	router.Use()
	// the jwt middleware
	authMiddleware, err := jwt.New(&middlewares.JwtMiddlewareStruct)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	publicV1 := router.Group("/api/v1")
	{
		publicV1.POST("/team", routers.CreateTeam)

		auth := publicV1.Group("/auth")
		auth.POST("/login", authMiddleware.LoginHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)

		services := publicV1.Group("/services")
		//services.GET("/teams/info", routers.CountTeamsHandler)
		services.GET("/system/info", routers.InfoHandler)
	}
	basicAuthV1 := router.Group("/api/v1")
	{
		// routes for prometheus with basic auth
		walker := basicAuthV1.Group("/game")
		walker.Use(gin.BasicAuth(gin.Accounts{"checker": config.Conf.CheckerPassword}))
		walker.GET("/checker", routers.CheckerHandler)
		walker.GET("/news", routers.NewsHandler)
		walker.GET("/exploit", routers.ExploitHandler)

	}
	jwtV1 := router.Group("/api/v1")
	jwtV1.Use(authMiddleware.MiddlewareFunc())
	{
		jwtV1.GET("/scoreboard", routers.ShowScoreboard)
		jwtV1.GET("/scoreboard/:name", routers.ShowTeamStatus)
		if config.Conf.Mode == "attack-defence" {
			jwtV1.POST("/submit", routers.SubmitFlagHandler)
		}

		// teams CRUD
		team := jwtV1.Group("/team")
		{
			team.GET("/", routers.GetTeamInfo)
			team.DELETE("/", routers.DeleteTeam)
		}

		// admins functions
		admin := jwtV1.Group("/admin")
		admin.Use(middlewares.IsAdmin())
		{
			admin.GET("/teams", routers.TeamsList)
			admin.POST("/vpn", routers.CreateVpnTeams)
			admin.DELETE("/team/:name", routers.DeleteTeams)
			admin.GET("/generate/variables", routers.GenerateVariables)
			admin.GET("/generate/sshkeys", routers.SshKeyArchiveHandler)
			admin.POST("/generate/prometheus", routers.GeneratePrometheus)
			admin.POST("/prom/start", routers.RunPrometheusHandler)
			admin.POST("/prom/stop", routers.StopPrometheusHandler)
			admin.GET("/reg/open", routers.OpenRegistrationHandler)
			admin.GET("/reg/close", routers.CloseRegistrationHandler)
		}
	}
	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Not found"})
	})
	router.Run(":8080")
}

package main

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/explabs/ad-ctf-paas-api/config"
	"github.com/explabs/ad-ctf-paas-api/database"
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

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

func isAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("id")
		if user.(*models.JWTTeam).TeamName == "admin" {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func AddAdmin() {
	teams, _ := database.GetUsers()
	for _, team := range teams {
		if team.Name == "admin" {
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
		Name:      "admin",
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

	database.InitMongo()

	var sc config.ServicesCost
	sc.Load()
	database.UploadServiceCost(sc.Services)

	database.InitRedis()

	AddAdmin()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.JWTTeam); ok {
				return jwt.MapClaims{
					identityKey: v.TeamName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.JWTTeam{
				TeamName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			filter := map[string]interface{}{
				"name": userID,
			}

			team, dbErr := database.FilterTeams(filter)
			if dbErr != nil {
				log.Println(dbErr)
				return nil, jwt.ErrFailedAuthentication
			}
			log.Println(team)
			if routers.CheckPasswordHash(password, team[0].Hash) {
				return &models.JWTTeam{
					TeamName: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.JWTTeam); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	// Simple group: v1
	v1 := router.Group("/api/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.POST("/submit", authMiddleware.MiddlewareFunc(), routers.SubmitFlagHandler)
		v1.GET("/scoreboard", authMiddleware.MiddlewareFunc(), routers.ShowScoreboard)
		v1.GET("/scoreboard/:name", authMiddleware.MiddlewareFunc(), routers.ShowTeamStatus)
		auth := router.Group("/auth")
		// Refresh time can be longer than token timeout
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
		auth.Use(authMiddleware.MiddlewareFunc())
		team := v1.Group("/team")
		team.POST("/", routers.CreateTeam)
		team.Use(authMiddleware.MiddlewareFunc())
		{
			team.GET("/", routers.GetTeamInfo)
			// team.POST("/")

			team.DELETE("/", routers.DeleteTeam)
		}
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.MiddlewareFunc())
		admin.Use(isAdmin())
		{
			admin.GET("/teams", routers.TeamsList)
			admin.POST("/vpn", routers.CreateVpnTeams)
			admin.DELETE("/team/:name", routers.DeleteTeams)
			// admin.POST("/generate/terraform", routers.GenerateTerraformConfig)
			admin.POST("/generate/variables", routers.GenerateVariables)
			admin.POST("/generate/sshkeys", routers.GenerateSshKeysDir)
			admin.POST("/generate/prometheus", routers.GeneratePrometheus)

		}
		services := v1.Group("/services")
		{
			services.GET("/teams/info", routers.CountTeamsHandler)
		}
		v1.GET("/checker",
			gin.BasicAuth(gin.Accounts{
				"checker": config.Conf.CheckerPassword,
			}),
			routers.CheckerHandler)
	}
	router.Run(":8080")
}

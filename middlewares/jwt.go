package middlewares

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/explabs/ad-ctf-paas-api/config"
	"github.com/explabs/ad-ctf-paas-api/database"
	"github.com/explabs/ad-ctf-paas-api/models"
	"github.com/explabs/ad-ctf-paas-api/routers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("id")
		if user.(*models.JWTTeam).TeamName == "admin" {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

var JwtMiddlewareStruct = jwt.GinJWTMiddleware{
	Realm:       "test zone",
	Key:         []byte("secret key"),
	Timeout:     time.Hour,
	MaxRefresh:  time.Hour,
	IdentityKey: identityKey,
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*models.JWTTeam); ok {
			return jwt.MapClaims{
				identityKey: v.TeamName,
				"mode":      config.Conf.Mode,
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
		team, dbErr := database.GetAuthTeam(userID)
		if dbErr != nil {
			log.Println(dbErr)
			return nil, jwt.ErrFailedAuthentication
		}
		log.Println(userID)
		if routers.CheckPasswordHash(password, team.Hash) {
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
	TokenLookup: "header: Authorization, query: token, cookie: jwt",
	TokenHeadName: "Bearer",
	TimeFunc: time.Now,
}

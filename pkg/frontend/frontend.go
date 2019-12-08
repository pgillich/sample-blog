// Package frontend is the frontend
package frontend

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pgillich/errfmt"

	"github.com/pgillich/sample-blog/api"
	"github.com/pgillich/sample-blog/internal/dao"
	"github.com/pgillich/sample-blog/internal/logger"
	"github.com/pgillich/sample-blog/internal/web"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

/*
func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*api.User).Name,
		"text":     "Hello World.",
	})
}
*/

// SetupGin is the service, called by automatic test, too
func SetupGin(router *gin.Engine, dbHandler *dao.Handler) *gin.Engine { //nolint:wsl
	authMiddleware, err := BuildAuthMiddleware(dbHandler)
	if err != nil {
		logger.Get().Panic("JWT Error, " + err.Error())
	}

	//nolint:gocritic
	/* TODO
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	*/

	v1 := router.Group("/api/v1")
	{ //nolint:gocritic
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.GET("/refresh_token", authMiddleware.RefreshHandler)

		v1.GET("/stat/user-post-comment", web.DecorHandlerDB(GetUserPostCommentStats, dbHandler))
	}

	return router
}

// BuildAuthMiddleware makes JWT middleware
func BuildAuthMiddleware(dbHandler *dao.Handler) (*jwt.GinJWTMiddleware, error) {
	identityKey := "id"

	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*api.User); ok {
				return jwt.MapClaims{
					identityKey: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &api.User{
				Name: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			userOK := dbHandler.CheckUser(loginVals.Username, loginVals.Password)
			// TODO admin user
			if userOK || (userID == "admin" && password == "admin") {
				return &api.User{
					Name: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// TODO admin user
			if v, ok := data.(*api.User); ok && v.Name == "admin" {
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

		// TimeFunc provides the current time. You can override it to use another time value.
		// This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}

// GetUserPostCommentStats collects and returns user activity
func GetUserPostCommentStats(c *gin.Context, dbHandler *dao.Handler) {
	if stats, err := dbHandler.GetUserPostCommentStats(c.Query("days")); err != nil {
		errs := logger.Get().WithError(err)
		statusCode := http.StatusBadRequest
		httpProblem := errfmt.BuildHTTPProblem(statusCode, errs)
		c.JSON(statusCode, httpProblem)

		errs.WithField("status", statusCode).Warning("cannot get stat")
	} else {
		c.JSON(http.StatusOK, stats)
	}
}

// Package frontend is the frontend
package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Depado/ginprom"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pgillich/errfmt"

	"github.com/pgillich/sample-blog/api"
	"github.com/pgillich/sample-blog/configs"
	"github.com/pgillich/sample-blog/internal/dao"
	"github.com/pgillich/sample-blog/internal/logger"
	"github.com/pgillich/sample-blog/internal/web"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// SetupGin is the service, called by automatic test, too
func SetupGin(router *gin.Engine, dbHandler *dao.Handler, enableMetrics bool) *gin.Engine { //nolint:wsl
	authMiddleware, err := BuildAuthMiddleware(dbHandler)
	if err != nil {
		logger.Get().Panic("JWT Error, " + err.Error())
	}

	//nolint:gocritic
	/* TODO
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	*/

	if enableMetrics {
		prom := ginprom.New(
			ginprom.Engine(router),
			ginprom.Subsystem("gin"),
			ginprom.Path("/metrics"),
		)
		router.Use(prom.Instrument())
	}

	v1 := router.Group("/api/v1")
	{ //nolint:gocritic
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.GET("/refresh_token", authMiddleware.RefreshHandler)

		v1.GET("/stat/user-post-comment", web.DecorHandlerDB(GetUserPostCommentStats, dbHandler))
	}
	v1.Use(authMiddleware.MiddlewareFunc())
	{ //nolint:gocritic
		v1.POST("/entry/:entry/comment", web.DecorHandlerDB(PostComment, dbHandler))
	}

	router.GET("/version", GetVersion)

	return router
}

// BuildAuthMiddleware makes JWT middleware
func BuildAuthMiddleware(dbHandler *dao.Handler) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: configs.AuthIdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*api.User); ok {
				return jwt.MapClaims{
					configs.AuthIdentityKey: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &api.User{
				Name: claims[configs.AuthIdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userName := loginVals.Username
			password := loginVals.Password

			userOK := dbHandler.CheckUser(loginVals.Username, loginVals.Password)
			// TODO admin user
			if userOK || (userName == "admin" && password == "admin") {
				return &api.User{
					Name: userName,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// TODO admin user
			if u, ok := data.(*api.User); ok /* && u.Name == "admin" */ {
				_, err := dbHandler.GetUserByName(u.Name)
				return err == nil
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

// PostComment creates a new comment to an entry
func PostComment(c *gin.Context, dbHandler *dao.Handler) {
	claims := jwt.ExtractClaims(c)
	userName := fmt.Sprintf("%s", claims[configs.AuthIdentityKey])
	id, _ := c.Get(configs.AuthIdentityKey)
	fmt.Print(id)

	var (
		err       error
		user      api.User
		entry     api.Entry
		entryID   int
		comment   api.Comment
		bodyBytes []byte
		body      api.Text
	)

	if bodyBytes, err = c.GetRawData(); err == nil {
		if err = json.Unmarshal(bodyBytes, &body); err == nil {
			if user, err = dbHandler.GetUserByName(userName); err == nil {
				if entryID, err = strconv.Atoi(c.Param("entry")); err == nil {
					if entry, err = dbHandler.GetEntryByID(uint(entryID)); err == nil {
						if user.ID != entry.UserID {
							err = fmt.Errorf("only own entry can be commented")
						} else {
							comment, err = dbHandler.CreateComment(user.ID, entry.ID, body.Text, time.Time{})
						}
					}
				}
			}
		}
	}

	// TODO refactor: move error handling and http status processing to common part
	if err != nil {
		errs := logger.Get().WithError(err)
		statusCode := http.StatusBadRequest
		httpProblem := errfmt.BuildHTTPProblem(statusCode, errs)
		c.JSON(statusCode, httpProblem)

		errs.WithField("status", statusCode).Warning("cannot get stat")
	} else {
		c.JSON(http.StatusOK, comment)
	}
}

// GetVersion returns build info
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, api.BuildInfo{
		Tag:       configs.BuildTag,
		Commit:    configs.BuildCommit,
		Branch:    configs.BuildBranch,
		BuildTime: configs.BuildTime,
	})
}

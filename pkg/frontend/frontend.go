// Package frontend is the frontend
package frontend

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pgillich/errfmt"

	"github.com/pgillich/sample-blog/internal/dao"
	"github.com/pgillich/sample-blog/internal/logger"
	"github.com/pgillich/sample-blog/internal/web"
)

// SetupGin is the service, called by automatic test, too
func SetupGin(router *gin.Engine, dbHandler *dao.Handler) *gin.Engine {
	//nolint:gocritic
	/*
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	*/

	v1 := router.Group("/api/v1")
	{ //nolint:gocritic
		//nolint:gocritic
		/*
			v1.Use(auth())
		*/
		v1.GET("/stat/user-post-comment", web.DecorHandlerDB(GetUserPostCommentStats, dbHandler))
	}

	return router
}

// GetUserPostCommentStats collects and returns user activity
func GetUserPostCommentStats(c *gin.Context, dbHandler *dao.Handler) {
	if stats, err := dbHandler.GetUserPostCommentStats(c.Param("days")); err != nil {
		errs := logger.Get().WithError(err)
		statusCode := http.StatusBadRequest
		httpProblem := errfmt.BuildHTTPProblem(statusCode, errs)
		c.JSON(statusCode, httpProblem)

		errs.WithField("status", statusCode).Warning("cannot get stat")
	} else {
		c.JSON(http.StatusOK, stats)
	}
}

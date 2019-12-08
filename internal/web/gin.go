// Package web provides gin related utils
package web

import (
	"github.com/gin-gonic/gin"

	"github.com/pgillich/sample-blog/internal/dao"
)

// HandlerDBFunc has DB parameter, too
type HandlerDBFunc func(c *gin.Context, dbHandler *dao.Handler)

// DecorHandlerDB decorates Gin handler with DB
func DecorHandlerDB(webHandlerFunc HandlerDBFunc, dbHandler *dao.Handler) func(c *gin.Context) {
	return func(c *gin.Context) {
		webHandlerFunc(c, dbHandler)
	}
}

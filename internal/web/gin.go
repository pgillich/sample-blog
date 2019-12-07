// Package web provides gin related utils
package web

import (
	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
)

// HandlerDBFunc has DB parameter, too
type HandlerDBFunc func(*gin.Context, *gorm.DB)

// DecorHandlerDB decorates Gin handler with DB
func DecorHandlerDB(dbHandler HandlerDBFunc, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dbHandler(c, db)
	}
}

// Package dao provides DB backends
package dao

import (
	"emperror.dev/errors"
	"github.com/jinzhu/gorm"

	"github.com/pgillich/sample-blog/api"
)

// Prepare inits tables.
func Prepare(db *gorm.DB) (*gorm.DB, error) {
	db.Exec("PRAGMA foreign_keys = ON")
	db.LogMode(true)

	db = db.AutoMigrate(&api.User{}, &api.Entry{}, &api.Comment{})
	if db.Error != nil {
		return nil, errors.Wrap(db.Error, "cannot update DB schema")
	}

	return db, nil
}

// GetUserPostCommentStats returns user activity stat
func GetUserPostCommentStats(db *gorm.DB) (api.UserPostCommentStats, error) {
	return api.UserPostCommentStats{}, nil
}

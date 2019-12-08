// Package dao provides DB backends
package dao

import (
	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // import Postgres driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // import Sqlite driver
	log "github.com/sirupsen/logrus"

	"github.com/pgillich/sample-blog/api"
	"github.com/pgillich/sample-blog/internal/logger"
)

// Handler is a thin layer over Gorm
type Handler struct {
	DB *gorm.DB
}

// NewHandler creates a new Gorm DB dbHandler
func NewHandler(dialect string, dsn string, sampleFill bool) (*Handler, error) {
	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to DB")
	}

	if logger.Get().IsLevelEnabled(log.DebugLevel) {
		db = db.Debug()
	}

	dbHandler := &Handler{DB: db}
	if err = dbHandler.prepare(); err != nil { //nolint:gocritic
		return nil, err
	}

	if sampleFill {
		if err = dbHandler.sampleFill(); err != nil { //nolint:gocritic
			return nil, err
		}
	}

	return dbHandler, nil
}

func (dbHandler *Handler) prepare() error {
	if dbHandler.DB.Dialect().GetName() == "sqlite3" {
		dbHandler.DB.Exec("PRAGMA foreign_keys = ON")
	}

	if logger.Get().IsLevelEnabled(log.DebugLevel) {
		dbHandler.DB = dbHandler.DB.LogMode(true)
	}

	db := dbHandler.DB.AutoMigrate(&api.User{}, &api.Entry{}, &api.Comment{})
	if db.Error != nil {
		return errors.Wrap(db.Error, "cannot update DB schema")
	}

	return nil
}

func (dbHandler *Handler) sampleFill() error {
	return nil
}

// Close closes the DB connection
func (dbHandler *Handler) Close() {
	dbHandler.DB.Close() //nolint:errcheck,gosec
}

// GetOrCreateUser creates a new user or returns, if exists
func (dbHandler *Handler) GetOrCreateUser(name string) (api.User, error) {
	templateUser := api.User{Name: name}
	user := api.User{}

	db := dbHandler.DB.Where(templateUser).FirstOrCreate(&user)
	if db.Error != nil {
		return templateUser, db.Error
	}

	return user, nil
}

// GetUserPostCommentStats returns user activity stat
func (dbHandler *Handler) GetUserPostCommentStats(days string) (api.UserPostCommentStats, error) {
	return api.UserPostCommentStats{}, nil
}

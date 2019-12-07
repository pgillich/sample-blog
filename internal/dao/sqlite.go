package dao

import (
	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // import Sqlite driver
	log "github.com/sirupsen/logrus"

	"github.com/pgillich/sample-blog/internal/logger"
)

// ConnectSqlite connects to the DB and inits it
func ConnectSqlite(dbDsn string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbDsn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to DB")
	}

	if logger.Get().IsLevelEnabled(log.DebugLevel) {
		db = db.Debug()
	}

	db, err = Prepare(db)

	return db, err
}

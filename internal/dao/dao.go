// Package dao provides DB backends
package dao

import (
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
	"strconv"
	"time"

	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // import Postgres driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // import Sqlite driver
	log "github.com/sirupsen/logrus"

	"github.com/pgillich/sample-blog/api"
	"github.com/pgillich/sample-blog/internal/logger"
)

// TimeNowFunc returns time.Now, hacked by automatic tests
type TimeNowFunc func() time.Time

// Handler is a thin layer over Gorm
type Handler struct {
	DB      *gorm.DB
	TimeNow TimeNowFunc
}

// NewHandler creates a new Gorm DB dbHandler
func NewHandler(dialect string, dsn string, samples []CompactSample) (*Handler, error) {
	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to DB")
	}

	if logger.Get().IsLevelEnabled(log.DebugLevel) {
		db = db.Debug()
	}

	dbHandler := &Handler{DB: db, TimeNow: time.Now}
	if err = dbHandler.prepare(); err != nil { //nolint:gocritic
		return nil, err
	}

	if len(samples) > 0 {
		if err = dbHandler.sampleFill(samples); err != nil { //nolint:gocritic
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

// EntryCommentSample is an entry --> comments map
type EntryCommentSample map[string][][]string

// CompactSample is a user --> entries --> comments structure
type CompactSample struct {
	api.User
	EntryComments []EntryCommentSample
}

// GetDefaultSampleFill returns the default samples for DB
func GetDefaultSampleFill() []CompactSample {
	return []CompactSample{
		{api.User{Name: "kovacsj", Password: "kovacs12"}, []EntryCommentSample{
			{"Entry K#1": [][]string{
				{"kovacsj", "Comment K#1/1", "2019-12-01"},
				{"szabop", "Comment K#1/2", "2019-12-01"},
				{"szabop", "Comment K#1/3", "2019-12-02"},
			}},
		}},
		{api.User{Name: "szabop", Password: "pal12"}, []EntryCommentSample{
			{"Entry S#1": [][]string{
				{"szabop", "Comment S#1/1", "2019-12-02"},
				{"kovacsj", "Comment S#1/2", "2019-12-02"},
				{"kovacsj", "Comment S#1/3", "2019-12-03"},
			}},
			{"Entry S#2": [][]string{
				{"kovacsj", "Comment S#2/1", "2019-12-03"},
			}},
			{"Entry S#3": [][]string{
				{"kovacsj", "Comment S#3/1", "2019-12-03"},
			}},
			{"Entry S#4": [][]string{
				{"kovacsj", "Comment S#4/1", "2019-12-04"},
				{"kovacsj", "Comment S#4/2", "2019-12-05"},
			}},
		}},
		{api.User{Name: "kocsisi", Password: "irma12"}, []EntryCommentSample{
			{"Entry I#1": [][]string{
				{"szabop", "Comment K#1/1", "2019-12-03"},
				{"kovacsj", "Comment I#1/2", "2019-12-03"},
			}},
		}},
	}
}

func (dbHandler *Handler) sampleFill(samples []CompactSample) error {
	users := map[string]api.User{}

	for _, userInfo := range samples {
		if user, err := dbHandler.CreateUser(userInfo.User.Name, userInfo.User.Password); err != nil {
			return err
		} else { //nolint:golint
			users[userInfo.User.Name] = user
		}
	}

	for _, userInfo := range samples {
		user := users[userInfo.User.Name]
		userID := user.ID

		for _, entryComments := range userInfo.EntryComments {
			for entryText, comments := range entryComments {
				entryTime, _ := time.Parse("2006-01-02", comments[0][2]) //nolint:errcheck

				if entry, err := dbHandler.CreateEntry(userID, entryText, entryTime); err != nil {
					return err
				} else { //nolint:golint
					for _, comment := range comments {
						commentUserID := users[comment[0]].ID
						commentText := comment[1]
						commentTime, _ := time.Parse("2006-01-02", comment[2]) //nolint:errcheck

						if _, err := dbHandler.CreateComment(commentUserID, entry.ID, commentText, commentTime); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// Close closes the DB connection
func (dbHandler *Handler) Close() {
	dbHandler.DB.Close() //nolint:errcheck,gosec
}

// CreateUser inserts a new user
func (dbHandler *Handler) CreateUser(name string, password string) (api.User, error) {
	user := api.User{Name: name, Password: hashPassword(password)}

	db := dbHandler.DB.Create(&user)
	if db.Error != nil {
		return user, errors.WrapWithDetails(db.Error, "cannot create user", "name", name)
	}

	return user, nil
}

// CreateEntry inserts a new blog entry
func (dbHandler *Handler) CreateEntry(userID uint, text string, ts time.Time) (api.Entry, error) {
	entry := api.Entry{UserID: userID, Text: text}
	if !ts.IsZero() {
		entry.CreatedAt = ts
		entry.UpdatedAt = ts
	}

	db := dbHandler.DB.Create(&entry)
	if db.Error != nil {
		return entry, errors.WrapWithDetails(db.Error, "cannot create entry", "userID", userID)
	}

	return entry, nil
}

// CreateComment inserts a new blog entry
func (dbHandler *Handler) CreateComment(userID uint, entryID uint, text string, ts time.Time) (api.Comment, error) {
	comment := api.Comment{UserID: userID, EntryID: entryID, Text: text}
	if !ts.IsZero() {
		comment.CreatedAt = ts
		comment.UpdatedAt = ts
	}

	db := dbHandler.DB.Create(&comment)
	if db.Error != nil {
		return comment, errors.WrapWithDetails(db.Error, "cannot create comment", "userID", userID, "entryID", entryID)
	}

	return comment, nil
}

// GetUserByName returns a user by name
func (dbHandler *Handler) GetUserByName(name string) (api.User, error) {
	var user api.User

	db := dbHandler.DB.Where(&api.User{Name: name}).First(&user)
	if db.Error != nil {
		return user, errors.Wrap(db.Error, "cannot get user")
	}

	return user, nil
}

// GetEntryByID returns a entry by id
func (dbHandler *Handler) GetEntryByID(id uint) (api.Entry, error) {
	var entry api.Entry

	db := dbHandler.DB.Where(&api.Entry{Model: gorm.Model{ID: id}}).First(&entry)
	if db.Error != nil {
		return entry, errors.Wrap(db.Error, "cannot get entry")
	}

	return entry, nil
}

// GetUserEntriesByName returns a entry by id
func (dbHandler *Handler) GetUserEntriesByName(name string) ([]api.Entry, error) {
	var entries []api.Entry

	user, err := dbHandler.GetUserByName(name)
	if err != nil {
		return entries, errors.Wrap(err, "cannot get user")
	}

	db := dbHandler.DB.Where(&api.Entry{UserID: user.ID}).Find(&entries)
	if db.Error != nil {
		return entries, errors.Wrap(db.Error, "cannot get entry")
	}

	return entries, nil
}

// CheckUser checks user for auth
func (dbHandler *Handler) CheckUser(name string, password string) bool {
	var user api.User

	user, err := dbHandler.GetUserByName(name)
	if err != nil {
		logger.Get().Info("User not exists", "name", name)
		return false
	}

	return user.Password == hashPassword(password)
}

// GetUserPostCommentStats returns user activity stat
// TODO refactoring: replace Unscoped().Raw() to better Gorm functions
func (dbHandler *Handler) GetUserPostCommentStats(daysString string) (api.UserPostCommentStats, error) {
	stats := api.UserPostCommentStats{}

	days, err := strconv.ParseUint(daysString, 10, 32)
	if err != nil {
		return stats, errors.WrapWithDetails(err, "invalid interval", "days", daysString)
	}

	now := dbHandler.TimeNow()
	fromTime := now.Add(-time.Duration(days*24) * time.Hour)

	users := []api.User{}
	db := dbHandler.DB.Find(&users)

	if db.Error != nil {
		return stats, errors.Wrap(db.Error, "cannot get users")
	}

	for _, user := range users {
		stats[user.ID] = &api.PostCommentStat{UserName: user.Name}
	}

	var userID, count uint

	//nolint:lll
	query := `SELECT user.id AS ID, COUNT(entry.id) AS Count FROM user JOIN entry ON entry.user_id = user.id WHERE entry.created_at >= ? GROUP BY user.id`
	rows, err := db.Unscoped().Raw(query, fromTime.Format(time.RFC3339)).Rows()

	if err != nil {
		return stats, errors.Wrap(err, "cannot get entries")
	}

	for rows.Next() {
		if err = rows.Scan(&userID, &count); err != nil {
			return stats, errors.Wrap(err, "cannot get entries")
		}

		stats[userID].Entries = count
	}

	//nolint:lll
	query = `SELECT user.id AS ID, COUNT(comment.id) AS Count FROM user LEFT JOIN comment ON comment.user_id = user.id WHERE comment.created_at >= ? GROUP BY user.id`

	rows, err = db.Unscoped().Raw(query, fromTime.Format(time.RFC3339)).Rows()
	if err != nil {
		return stats, errors.Wrap(err, "cannot get entries")
	}

	for rows.Next() {
		if err := rows.Scan(&userID, &count); err != nil {
			return stats, errors.Wrap(err, "cannot get entries")
		}

		stats[userID].Comments = count
	}

	return stats, nil
}

func hashPassword(password string) string {
	h := sha1.New()           //nolint:gosec
	h.Write([]byte(password)) //nolint:errcheck,gosec

	return hex.EncodeToString(h.Sum(nil))
}

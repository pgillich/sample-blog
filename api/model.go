/* Package api provides external type definitions
It should not import any local packages.
*/
package api

import (
	"github.com/jinzhu/gorm"
)

// User is the internal representation of user table
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"uniq_key"`
	Password string `json:"-"`
}

// TableName forces table name singular
func (User) TableName() string {
	return "user"
}

// Entry is the internal representation of entry table
type Entry struct {
	gorm.Model
	UserID uint   `json:"userID" sql:"type:integer REFERENCES user(id)"`
	Text   string `json:"text"`
}

// TableName forces table name singular
func (Entry) TableName() string {
	return "entry"
}

// Comment is the internal representation of comment table
type Comment struct {
	gorm.Model
	UserID  uint   `json:"userID" sql:"type:integer REFERENCES user(id)"`
	EntryID uint   `json:"entryID" sql:"type:integer REFERENCES entry(id)"`
	Text    string `json:"text"`
}

// TableName forces table name singular
func (Comment) TableName() string {
	return "comment"
}

// PostCommentStat is statistic about a user activity
type PostCommentStat struct {
	UserName string `json:"userName"`
	Entries  uint   `json:"entries"`
	Comments uint   `json:"comments"`
}

// UserPostCommentStats is statistic about user activities
type UserPostCommentStats map[uint]*PostCommentStat

// Text is a simple text struct
type Text struct {
	Text string `json:"text"`
}

// BuildInfo is the schema to /version
type BuildInfo struct {
	Tag       string `json:"tag"`
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
	BuildTime string `json:"buildTime"`
}

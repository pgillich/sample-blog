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

// Entry is the internal representation of entry table
type Entry struct {
	gorm.Model
	UserID uint   `json:"userID" sql:"type:integer REFERENCES users(id)"`
	Text   string `json:"text"`
}

// Comment is the internal representation of comment table
type Comment struct {
	gorm.Model
	UserID  uint   `json:"userID" sql:"type:integer REFERENCES users(id)"`
	EntryID uint   `json:"entryID" sql:"type:integer REFERENCES entries(id)"`
	Text    string `json:"text"`
}

// PostCommentStat is statistic about a user activity
type PostCommentStat struct {
	UserName    string `json:"userName"`
	EntryText   uint   `json:"entryText"`
	CommentText uint   `json:"commentText"`
}

// UserPostCommentStats is statistic about user activities
type UserPostCommentStats map[uint]PostCommentStat

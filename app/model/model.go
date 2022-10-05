package model

import (
	"blog/config"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var dbc = config.Cfg.Database

func init() {

	DSL := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=%s&loc=%s",
		dbc.User, dbc.Password, dbc.Host, dbc.Name, dbc.Charset, dbc.ParseTime, dbc.Loc)

	database, err := gorm.Open(dbc.Dialect, DSL)

	if err != nil {
		panic("Failed to connect to database!")
	}
	database.LogMode(true)
	DB = database
}

type User struct {
	gorm.Model
	UserName string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Passwd   string
	Github   string
	Articles []Article `gorm:"foreignKey:UserID;references:UserID"`
	Comments []Comment `gorm:"foreignKey:CommentID;references:CommentID"`
}

// 设置 `User` 的表名为 `profiles`
func (User) TableName() string {
	return "Users"
}

type Article struct {
	gorm.Model
	Title    string `gorm:"unique;not null;varchar(100)"`
	Content  string `gorm:"unique;not null;"`
	Tags     string `gorm:"varchar(100)"`
	UserID   uint
	Comments []Comment `gorm:"foreignKey:CommentID;references:CommentID"`
}

type Comment struct {
	gorm.Model
	Content   string `gorm:"unique;not null;"`
	UserID    uint
	ArticleID uint
}

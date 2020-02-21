package models

import (
	"time"

	log "github.com/halaalajlan/hackathon/logger"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Response contains the attributes found in an API response
type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username" sql:"not null;unique"`
	Hash     string `json:"-"`
	ApiKey   string `json:"api_key" sql:"not null;unique"`
}

func SetUp() error {

	// Open our database connection
	i := 0
	var err error
	for {
		db, err = gorm.Open("sqlite3", "hackathon.db")
		if err == nil {
			break
		}
		if err != nil && i >= 10 {
			log.Error(err)
			return err
		}
		i += 1
		log.Warn("waiting for database to be up...")
		time.Sleep(5 * time.Second)
	}
	db.LogMode(false)
	db.SetLogger(log.Logger)
	db.DB().SetMaxOpenConns(1)
	return nil

}

// GetUserByUsername returns the user that the given username corresponds to. If no user is found, an
// error is thrown.
func GetUserByUsername(username string) (User, error) {
	u := User{}
	err := db.Where("username = ?", username).First(&u).Error
	return u, err
}

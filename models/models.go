package models

import (
	"time"

	log "github.com/halaalajlan/hackathon/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // Blank import needed to import sqlite3
)

var db *gorm.DB

// Response contains the attributes found in an API response
type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// Flash is used to hold flash information for use in templates.
type Flash struct {
	Type    string
	Message string
}

type Hospital struct {
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
		log.Error(err)
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
func GetUserByUsername(username string) (Hospital, error) {
	u := Hospital{}
	err := db.Table("Hospital").Where("name_hospital = ?", username).First(&u).Error
	return u, err
}

// GetUserByUsername returns the user that the given username corresponds to. If no user is found, an
// error is thrown.
func GetUserByAPIKey(username string) (Hospital, error) {
	u := Hospital{}
	err := db.Table("Hospital").Where("api_Key = ?", username).First(&u).Error
	return u, err
}

// GetUser returns the user that the given id corresponds to. If no user is found, an
// error is thrown.
func GetUser(id int64) (Hospital, error) {
	u := Hospital{}
	err := db.Table("Hospital").Where("ID_Hosptial=?", id).First(&u).Error
	return u, err
}

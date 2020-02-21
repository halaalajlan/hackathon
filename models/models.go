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
	Id            int64 `gorm:"column:id; primary_key:yes"`
	name_hospital string
	api_Key       int64
	city          string
	Address_H     string
	phone         int64
	Email_Admin   string
	Hash          string
	ApiKey        string `json:"api_key" sql:"not null;unique"`
}

type Medical_Record struct {
	Id                  int64  ` gorm:"column:id; primary_key:yes"`
	Diabetic            string `gorm:"column:diabetic`
	High_Blood_Pressure string `gorm:"column:high_blood_pressure`
	Cholestrol          string `gorm:"column:cholestrol`
	Heart_dieases       string `gorm:"column:heart_dieases`
	Asthma              string `gorm:"column:asthma`
	Allergic_disease    string `gorm:"column:allergic_disease`
	Id_Patent           int    `gorm:"column:id_patent`
}

type Patient struct {
	Id           int64 `json:"id" gorm:"column:id; primary_key:yes"`
	Fname        string
	Lname        string
	Birthday     time.Time
	typeOfBlood  string
	api_Key      int64
	gander       string
	Phone_number int64
	home_address string
}

type Pat_Hospital struct {
	Id_Patent       int64
	ID_Hosptial     int64
	Last_Visit_Date time.Time
	ReasonForVisit  string
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

func GetPatientRecord(id int) (Medical_Record, error) {
	md := Medical_Record{}
	err := db.Table("Medical_Record").Where("Id_Patent=?", id).Find(&md).Error
	return md, err
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
	err := db.Table("Hospital").Where("id=?", id).First(&u).Error
	return u, err
}

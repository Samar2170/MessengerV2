package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open(postgres.Open(DBURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Db.AutoMigrate(&User{})
	Db.AutoMigrate(&Service{})
	Db.AutoMigrate(&Subscriber{})
	Db.AutoMigrate(&Subscriptions{})

}

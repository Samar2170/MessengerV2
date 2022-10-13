package main

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique_index"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique_index"`
}
type Service struct {
	*gorm.Model
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"unique_index"`
	User   User   `gorm:"foreignkey:UserID"`
	UserID uint
}
type Subscriber struct {
	*gorm.Model
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"unique_index"`
	ChatID    int64  `gorm:"unique_index"`
	FirstName string
	LastName  string
}
type Subscriptions struct {
	*gorm.Model
	ID           uint       `gorm:"primary_key"`
	Subscriber   Subscriber `gorm:"foreignkey:SubscriberID"`
	SubscriberID uint
	Service      Service `gorm:"foreignkey:ServiceID"`
	ServiceID    uint
}

func (u User) Create() error {
	err := Db.Create(&u).Error
	return err
}

func FindUser(username string) (User, error) {
	var u User
	err := Db.First(&u, "username = ?", username).Error
	return u, err
}
func GetUserId(username string) (uint, error) {
	var u User
	err := Db.First(&u, "username = ?", username).Error
	return u.ID, err
}

func CreateService(uid uint, serviceName string) error {
	s := Service{Name: serviceName, UserID: uid}
	err := Db.Create(&s).Error
	return err
}
func GetService(serviceName string) (Service, error) {
	var s Service
	err := Db.First(&s, "name = ?", serviceName).Error
	return s, err
}

func (s Subscriber) Create() error {
	err := Db.FirstOrCreate(&s, "chat_id = ?", s.ChatID).Error
	return err
}
func GetSubscriber(chatID int64) (Subscriber, error) {
	var s Subscriber
	err := Db.First(&s, "chat_id = ?", chatID).Error
	return s, err
}
func (s Subscriptions) Create() error {
	err := Db.Create(&s).Error
	return err
}

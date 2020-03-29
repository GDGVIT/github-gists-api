package user

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name       string `json:"Name"`
	Email      string `json:"Email"`
	OAuthToken string `json:"-"`
}

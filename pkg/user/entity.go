package user

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name       string `json:"name"`
	Email      string `json:"email"`
	OAuthToken string `json:"oauth_token"`
}

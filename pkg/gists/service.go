package gists

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type Service interface {
	ViewAllFiles(userID float64) (*[]Gist, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}

var client = &http.Client{Timeout: 10 * time.Second}

func (s *service) ViewAllFiles(userID float64) (*[]Gist, error) {
	/*
		user := &u.User{}
		err := s.db.Where("id=?", userID).First(user)
		if err != nil {
			return nil, pkg.ErrDatabase
		}
	*/

	// Accessing github for the files
	//token := user.OAuthToken
	req, er := http.NewRequest("GET", "https://api.github.com/gists", nil)
	if er != nil {
		return nil, er
	}
	req.Header.Set("Authorization", "token 9889520a036bc1d19744ecae4667c1f375ab4d88")
	resp, er := client.Do(req)

	if er != nil {
		return nil, er
	}

	var gists []Gist
	er = json.NewDecoder(resp.Body).Decode(&gists)
	if er != nil {
		return nil, er
	}
	return &gists, nil
}

package gists

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/GistsBackend/pkg"
	"github.com/rithikjain/GistsBackend/pkg/user"
	"io/ioutil"
	"net/http"
	"time"
)

type Service interface {
	ViewAllFiles(userID float64) (*[]File, error)
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

func (s *service) ViewAllFiles(userID float64) (*[]File, error) {
	user := &user.User{}
	err := s.db.Where("id=?", userID).First(user).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	// Accessing github for the files
	token := user.OAuthToken
	req, er := http.NewRequest("GET", "https://api.github.com/gists", nil)
	if er != nil {
		return nil, er
	}
	req.Header.Set("Authorization", "token "+token)
	resp, er := client.Do(req)

	if er != nil {
		return nil, er
	}

	var gists []Gist
	er = json.NewDecoder(resp.Body).Decode(&gists)
	if er != nil {
		return nil, er
	}

	var files []File

	// Logic for formatting the files in right format
	for i := 0; i < len(gists); i++ {
		for _, file := range gists[i].Files {
			file.GistID = gists[i].ID
			file.GistUrl = gists[i].Url
			file.IsPublic = gists[i].IsPublic
			file.UpdatedAt = gists[i].UpdatedAt
			file.Description = gists[i].Description

			// Get content of the file
			res, err := http.Get(file.RawUrl)
			if err != nil {
				return nil, err
			}
			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			content := string(resBody)
			file.Content = content

			files = append(files, file)
		}
	}

	return &files, nil
}

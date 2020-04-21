package gists

import (
	"bytes"
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

	CreateGist(userID float64, gistFile *GistFile) (*[]File, error)

	UpdateGist(userID float64, gistFile *GistFile) (*[]File, error)

	DeleteGist(userID float64, deleteGist *DeleteGist) (*[]File, error)
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
	defer resp.Body.Close()
	return &files, nil
}

func (s *service) CreateGist(userID float64, gistFile *GistFile) (*[]File, error) {
	// Transforming the file into correct format to send to the gist api
	fileMap := make(map[string]FileContent)
	fileMap[gistFile.Filename] = FileContent{Content: gistFile.Content}
	request := CreateFileRequest{
		Description: gistFile.Description,
		IsPublic:    gistFile.IsPublic,
		Files:       fileMap,
	}

	requestJson, _ := json.Marshal(request)

	user := &user.User{}
	err := s.db.Where("id=?", userID).First(user).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	// Creating file on github
	token := user.OAuthToken
	req, er := http.NewRequest("POST", "https://api.github.com/gists", bytes.NewBuffer(requestJson))
	if er != nil {
		return nil, er
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, er := client.Do(req)
	if er != nil {
		return nil, er
	}

	var gist Gist
	er = json.NewDecoder(resp.Body).Decode(&gist)
	if er != nil {
		return nil, er
	}
	var files []File
	for _, file := range gist.Files {
		file.GistID = gist.ID
		file.GistUrl = gist.Url
		file.IsPublic = gist.IsPublic
		file.UpdatedAt = gist.UpdatedAt
		file.Description = gist.Description

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
	defer resp.Body.Close()
	return &files, nil
}

func (s *service) UpdateGist(userID float64, gistFile *GistFile) (*[]File, error) {
	// Transforming the file into correct format to send to the gist api
	fileMap := make(map[string]FileContent)
	fileMap[gistFile.Filename] = FileContent{Content: gistFile.Content}
	request := UpdateFileRequest{
		Description: gistFile.Description,
		Files:       fileMap,
	}

	requestJson, _ := json.Marshal(request)

	user := &user.User{}
	err := s.db.Where("id=?", userID).First(user).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	// Updating the gist on github
	token := user.OAuthToken
	req, er := http.NewRequest("PATCH", "https://api.github.com/gists/"+gistFile.GistID, bytes.NewBuffer(requestJson))
	if er != nil {
		return nil, er
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, er := client.Do(req)
	if er != nil {
		return nil, er
	}

	var gist Gist
	er = json.NewDecoder(resp.Body).Decode(&gist)
	if er != nil {
		return nil, er
	}
	var files []File
	for _, file := range gist.Files {
		file.GistID = gist.ID
		file.GistUrl = gist.Url
		file.IsPublic = gist.IsPublic
		file.UpdatedAt = gist.UpdatedAt
		file.Description = gist.Description

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
	defer resp.Body.Close()
	return &files, nil
}

func (s *service) DeleteGist(userID float64, deleteGist *DeleteGist) (*[]File, error) {
	// Transforming the file into correct format to send to the gist api
	file := make(map[string]interface{})
	file[deleteGist.Filename] = nil

	request := make(map[string]interface{})
	request["files"] = file

	requestJson, _ := json.Marshal(request)

	user := &user.User{}
	err := s.db.Where("id=?", userID).First(user).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	// Deleting file from github
	token := user.OAuthToken
	req, er := http.NewRequest("PATCH", "https://api.github.com/gists/"+deleteGist.GistID, bytes.NewBuffer(requestJson))
	if er != nil {
		return nil, er
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, er := client.Do(req)
	if er != nil {
		return nil, er
	}

	if resp.Header.Get("Status") == "404 Not Found" {
		return nil, pkg.ErrNotFound
	}

	var gist Gist
	er = json.NewDecoder(resp.Body).Decode(&gist)
	if er != nil {
		return nil, er
	}
	var files []File
	for _, file := range gist.Files {
		file.GistID = gist.ID
		file.GistUrl = gist.Url
		file.IsPublic = gist.IsPublic
		file.UpdatedAt = gist.UpdatedAt
		file.Description = gist.Description

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
	defer resp.Body.Close()
	return &files, nil
}

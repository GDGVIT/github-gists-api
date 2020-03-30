package handler

import (
	"encoding/json"
	"github.com/rithikjain/GistsBackend/api/middleware"
	"github.com/rithikjain/GistsBackend/api/view"
	"github.com/rithikjain/GistsBackend/pkg/gists"
	"io/ioutil"
	"net/http"
)

// Protected Request
func viewAllFiles(s gists.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}

		g, err := s.ViewAllFiles(claims["id"].(float64))
		if err != nil {
			view.Wrap(err, w)
			return
		}

		var files []gists.File
		var allGists []gists.Gist
		mGists := *g

		// Logic for formatting the files in right format
		for i := 0; i < len(mGists); i++ {
			allGists = append(allGists, mGists[i])
			for _, file := range mGists[i].Files {
				file.GistID = allGists[i].ID
				file.GistUrl = allGists[i].Url
				file.IsPublic = allGists[i].IsPublic
				file.UpdatedAt = allGists[i].UpdatedAt
				file.Description = allGists[i].Description

				// Get content of the file
				res, err := http.Get(file.RawUrl)
				if err != nil {
					view.Wrap(err, w)
					return
				}
				resBody, err := ioutil.ReadAll(res.Body)
				if err != nil {
					view.Wrap(err, w)
					return
				}
				content := string(resBody)
				file.Content = content

				files = append(files, file)
			}
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Files Retrieved",
			"status":  http.StatusOK,
			"files":   files,
		})
	})
}

func MakGistsHandler(r *http.ServeMux, svc gists.Service) {
	r.Handle("/api/gists", middleware.Validate(viewAllFiles(svc)))
}

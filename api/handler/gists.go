package handler

import (
	"encoding/json"
	"github.com/rithikjain/GistsBackend/api/middleware"
	"github.com/rithikjain/GistsBackend/api/view"
	"github.com/rithikjain/GistsBackend/pkg/gists"
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

		files, err := s.ViewAllFiles(claims["id"].(float64))
		if err != nil {
			view.Wrap(err, w)
			return
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

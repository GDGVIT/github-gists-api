package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/rithikjain/GistsBackend/api/middleware"
	"github.com/rithikjain/GistsBackend/api/view"
	"github.com/rithikjain/GistsBackend/pkg/user"
	"net/http"
	"os"
)

func register(s user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
		}

		user := &user.User{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			view.Wrap(err, w)
			return
		}

		u, err := s.Register(user)
		if err != nil {
			view.Wrap(err, w)
			return
		}

		// Handling JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   u.ID,
			"role": "user",
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("jwt_secret")))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		u.OAuthToken = ""

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Account Registered",
			"status":  http.StatusCreated,
			"token":   tokenString,
			"user":    u,
		})
	})
}

// Protected request
func userDetails(s user.Service) http.Handler {
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
		u, err := s.GetUserByID(claims["id"].(float64))
		if err != nil {
			view.Wrap(err, w)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User Found",
			"status":  http.StatusOK,
			"user":    u,
		})
	})
}

func MakeUserHandler(r *http.ServeMux, svc user.Service) {
	r.Handle("/api/user/register", register(svc))
	r.Handle("/api/user/details", middleware.Validate(userDetails(svc)))
}

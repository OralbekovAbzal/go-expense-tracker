package authorization

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Users []User

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if http.MethodPost == r.Method {
		var newUser User

		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Invalid login or password", http.StatusBadRequest)
			return
		}

		for _, user := range Users {
			if user.Username == newUser.Username {
				http.Error(w, "Username already exists", http.StatusConflict)
				return
			}
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		newUser.Password = string(hash)
		Users = append(Users, newUser)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if http.MethodPost == r.Method {
		var userLogin User

		err := json.NewDecoder(r.Body).Decode(&userLogin)
		if err != nil {
			http.Error(w, "Invalid login or password", http.StatusBadRequest)
			return
		}

		for _, user := range Users {
			if user.Username == userLogin.Username {
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
				if err == nil {
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
						"username": user.Username,
						"exp":      time.Now().Add(time.Hour * 24).Unix(),
					})

					tokenString, err := token.SignedString([]byte("super-secret-key"))
					if err != nil {
						http.Error(w, "Failed to generate token", http.StatusInternalServerError)
						return
					}

					json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
					return
				}
			}
		}
	}
}

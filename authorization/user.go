package authorization

import (
	db2 "API_service/db"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type User struct {
	Id       int    `json:"id"`
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

		db := db2.ConnectDataBase()
		query := `insert into users (username,password) values ($1,$2)`

		_, err = db.Exec(query, newUser.Username, newUser.Password)
		db.Close()
		if err != nil {
			http.Error(w, "Error creating new user", http.StatusInternalServerError)
			fmt.Println("DB error:", err)
			return
		}

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

		db := db2.ConnectDataBase()

		query := `Select * from users where username=$1`
		row := db.QueryRow(query, userLogin.Username)

		var id int
		var username, password string

		err = row.Scan(&id, &username, &password)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Println("User not found")
			} else {
				fmt.Println("Query error", err)
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(userLogin.Password))
		if err == nil {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id": id,
				"exp":     time.Now().Add(time.Hour * 24).Unix(),
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

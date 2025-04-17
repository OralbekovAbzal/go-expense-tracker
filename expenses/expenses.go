package expenses

import (
	db2 "API_service/db"
	"API_service/middleware"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

type Expense struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
	UserId int     `json:"user_id"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{"message": "Hello, Abzal"}
	json.NewEncoder(w).Encode(response)
}

func AddExpenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID := r.Context().Value(middleware.User_idKey).(int)

		var newExpense Expense

		err := json.NewDecoder(r.Body).Decode(&newExpense)
		if err != nil {
			fmt.Println("Invalid data")
			return
		}

		newExpense.UserId = userID

		db := db2.ConnectDataBase()
		query := `Insert into expenses (title,amount,user_id) values ($1,$2,$3)`

		_, err = db.Exec(query, newExpense.Title, newExpense.Amount, newExpense.UserId)
		if err != nil {
			http.Error(w, "Error creating new expense", http.StatusInternalServerError)
			fmt.Println("DB error:", err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Expense added"})
	}
}

func AllExpensesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.User_idKey).(int)

	db := db2.ConnectDataBase()
	query := `select * from expenses where user_id=$1`

	rows, err := db.Query(query, userID)
	if err != nil {
		http.Error(w, "Error searching", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var expenses []Expense

	for rows.Next() {
		var expense Expense

		err = rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.UserId)
		if err != nil {
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			fmt.Println("Scan error:", err)
			return
		}

		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		fmt.Println("Iteration error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func deleteExpense(w http.ResponseWriter, r *http.Request) {
	if http.MethodDelete == r.Method {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid expense id", http.StatusBadRequest)
			return
		}

		query := `Delete from expenses where id=$1`

		db := db2.ConnectDataBase()

		_, err = db.Exec(query, id)
		if err != nil {
			http.Error(w, "Error quering database", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Expense deleted")
	}
}

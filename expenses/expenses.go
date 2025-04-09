package expenses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Expense struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
}

var expenses []Expense
var nextID = 1

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{"message": "Hello, Abzal"}
	json.NewEncoder(w).Encode(response)
}

func AddExpenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newExpense Expense

		err := json.NewDecoder(r.Body).Decode(&newExpense)
		if err != nil {
			fmt.Println("Invalid data")
			return
		}
		newExpense.ID = nextID
		nextID++
		expenses = append(expenses, newExpense)
	}
}

func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

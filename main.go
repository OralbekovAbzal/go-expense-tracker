package main

import (
	"API_service/authorization"
	"API_service/expenses"
	"API_service/middleware"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", expenses.HelloHandler)
	http.HandleFunc("/add", middleware.AuthMiddleware(expenses.AddExpenseHandler))
	http.HandleFunc("/all", middleware.AuthMiddleware(expenses.AllExpensesHandler))

	http.HandleFunc("/register", authorization.RegisterHandler)
	http.HandleFunc("/login", authorization.LoginHandler)
	http.HandleFunc("/me", middleware.AuthMiddleware(expenses.HelloHandler))

	fmt.Println("Server is starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

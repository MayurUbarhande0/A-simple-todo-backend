package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MayurUbarhande0/backend/db"
	"github.com/MayurUbarhande0/backend/internals/auth"
	"github.com/MayurUbarhande0/backend/internals/middleware"
	"github.com/MayurUbarhande0/backend/internals/todo"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	conn := db.Connect()
	defer conn.Close()

	authHandler := auth.NewAuthHandler(conn)
	todo.SetDB(conn)

	r := mux.NewRouter()

	r.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/todos", todo.GetAllTasks).Methods("GET")
	protected.HandleFunc("/todos", todo.CreateTask).Methods("POST")
	protected.HandleFunc("/todos/{id}", todo.UpdateTask).Methods("PUT")
	protected.HandleFunc("/todos/{id}", todo.DeleteTask).Methods("DELETE")

	port := os.Getenv("PORT")
	fmt.Println("server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

package main

import (
	"TestApi/Services"
	_ "TestApi/docs"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

// @title Test API
// @version 1.0
// @description This is a sample API for testing Swagger integration
// @host localhost:5555
// @BasePath /
func main() {
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)

	router := mux.NewRouter()
	if err != nil {
		log.Printf("Error while connecting to DB: %v", err)
	}
	_, _ = db.Exec("create table if not exists users (name varchar(255) not null, surname varchar(255), email   varchar(255) not null, id   serial   constraint users_pk   primary key );  ")

	router.HandleFunc("/user/all", Services.GetAllUsers).Methods("GET")
	router.HandleFunc("/user/getById/{id}", Services.GetUserById).Methods("GET")
	router.HandleFunc("/user/create", Services.CreateUser).Methods("POST")
	router.HandleFunc("/user/delete/{id}", Services.RemoveUserById).Methods("DELETE")
	router.HandleFunc("/user/update", Services.UpdateUserById).Methods("PATCH")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	_ = http.ListenAndServe(":8080", router)
}

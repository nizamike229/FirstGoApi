package Services

import (
	"TestApi/Models"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/ahmetb/go-linq/v3"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
	_ "github.com/samber/lo/parallel"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserRegisterModel struct {
	Name    string
	Surname string
	Email   string
}

var db *sql.DB
var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"))

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var Users []Models.User
	rows, err := db.Query("SELECT id as Id,name as Name,surname as Surname,email as Email FROM users")
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		u := Models.User{}
		err := rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Email)
		if err != nil {
			panic(err)
		}
		Users = append(Users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(Users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	Users := getAllUserPrivate()
	inputId := mux.Vars(r)["id"]
	filteredUsers := lo.Filter(Users, func(item Models.User, index int) bool {
		return strconv.Itoa(item.Id) == inputId
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(filteredUsers)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var userRequest UserRegisterModel
	_ = json.NewDecoder(r.Body).Decode(&userRequest)

	newUser := Models.User{
		Name:    userRequest.Name,
		Surname: userRequest.Surname,
		Email:   userRequest.Email,
	}

	_, err := db.Exec("insert into Users (name, surname, email) values ($1, $2, $3)", newUser.Name, newUser.Surname, newUser.Email)
	if err != nil {
		log.Printf("Error while creating user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("User with email:  %s added successfully.", newUser.Email)

	fmt.Println(message)
	_ = json.NewEncoder(w).Encode(message)
}

func RemoveUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("delete from Users where id=$1", id)
	if err != nil {
		log.Printf("Error while removing user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode("User was removed successfully!")
}

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var editRequest Models.User
	_ = json.NewDecoder(r.Body).Decode(&editRequest)

	users := getAllUserPrivate()
	isExist := lo.ContainsBy(users, func(u Models.User) bool {
		return u.Id == editRequest.Id
	})
	if !isExist {
		_ = json.NewEncoder(w).Encode("User not found")
		return
	}

	if editRequest.Name != "" {
		_, err := db.Exec("update Users set Name = $1 where Id = $2", editRequest.Name, editRequest.Id)
		if err != nil {
			log.Printf("Error while updating user data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	if editRequest.Surname != "" {
		_, err := db.Exec("update Users set Surname = $1 where Id = $2", editRequest.Surname, editRequest.Id)
		if err != nil {
			log.Printf("Error while updating user data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	if editRequest.Email != "" {
		_, err := db.Exec("update Users set Email = $1 where Id = $2", editRequest.Email, editRequest.Id)
		if err != nil {
			log.Printf("Error while updating user data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	_ = json.NewEncoder(w).Encode("User was updated successfully!")
}

func getAllUserPrivate() []Models.User {
	var Users []Models.User
	rows, err := db.Query("SELECT id as Id,name as Name,surname as Surname,email as Email FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := Models.User{}
		err := rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Email)
		if err != nil {
			panic(err)
		}
		Users = append(Users, u)
	}

	return Users
}

func init() {
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf(err.Error())
		}

		if err = db.Ping(); err != nil {
			log.Printf(err.Error())
		}
		time.Sleep(2 * time.Second)
	}

}

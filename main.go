package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

var db *gorm.DB

func init() {
	var err error
	dsn := "user=postgres dbname=postgres password=admin sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Миграция таблицы User
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
}

// Создание нового пользователя
func createUser(user *User) error {
	return db.Create(user).Error
}

// Получение пользователя по ID
func getUserByID(id uint) (User, error) {
	var user User
	err := db.First(&user, id).Error
	return user, err
}

// Обновление имени пользователя по ID
func updateUsernameByID(id uint, newUsername string) error {
	return db.Model(&User{}).Where("id = ?", id).Update("username", newUsername).Error
}

// Удаление пользователя по ID
func deleteUserByID(id uint) error {
	return db.Delete(&User{}, id).Error
}

// Получение списка всех пользователей
func getAllUsers() ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

type JsonRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Message  string `json:"message"`
}

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	//newUser := User{
	//	Username: "test2",
	//	Email:    "asd@asd",
	//	Password: "qwe",
	//}
	//// Создание нового пользователя
	//err := createUser(&newUser)
	//if err != nil {
	//	fmt.Println("Error creating user:", err)
	//	return
	//}
	//fmt.Println("User created successfully. ID:", newUser.ID)
	//
	//// Получение пользователя по ID
	//userByID, err := getUserByID(newUser.ID)
	//if err != nil {
	//	fmt.Println("Error getting user by ID:", err)
	//	return
	//}
	//fmt.Printf("User by ID %d: %+v\n", newUser.ID, userByID)
	//
	//// Обновление имени пользователя по ID
	//err = updateUsernameByID(newUser.ID, "new_john_doe")
	//if err != nil {
	//	fmt.Println("Error updating username:", err)
	//	return
	//}
	//fmt.Println("Username updated successfully.")
	//
	//// Получение списка всех пользователей
	//allUsers, err := getAllUsers()
	//if err != nil {
	//	fmt.Println("Error getting all users:", err)
	//	return
	//}
	//fmt.Println("All Users:")
	//for _, u := range allUsers {
	//	fmt.Printf("ID: %d, Username: %s, Email: %s\n", u.ID, u.Username, u.Email)
	//}
	//
	//// Удаление пользователя по ID
	//err = deleteUserByID(newUser.ID)
	//if err != nil {
	//	fmt.Println("Error deleting user:", err)
	//	return
	//}
	//fmt.Println("User deleted successfully.")
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Добавляем поддержку CORS для запросов OPTIONS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodPost:
		handlePostRequest(w, r)
	case http.MethodGet:
		handleGetRequest(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	var requestData JsonRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, `{"status": "400", "message": "Invalid JSON message"}`, http.StatusBadRequest)
		return
	}

	if requestData.Message == "" && (requestData.Username == "" || requestData.Email == "" || requestData.Password == "") {
		http.Error(w, `{"status": "400", "message": "Either 'message' or all 'username', 'email', 'password' fields must be provided"}`, http.StatusBadRequest)
		return
	}

	if requestData.Message != "" {
		fmt.Println("Received POST message:", requestData.Message)
	} else {
		fmt.Printf("Received POST request:\nUsername: %s\nEmail: %s\nPassword: %s\n", requestData.Username, requestData.Email, requestData.Password)
	}
	//saveRegistrationData(requestData.Username, requestData.Email, requestData.Password)

	// Добавляем заголовки CORS для основного запроса
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := JsonResponse{
		Status:  "success",
		Message: "Data successfully received",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	name := queryParams.Get("name")
	age := queryParams.Get("age")

	if name == "" || age == "" {
		http.Error(w, `{"status": "400", "message": "Both name and age parameters are required in the GET request"}`, http.StatusBadRequest)
		return
	}

	fmt.Printf("Received GET request with name: %s, age: %s\n", name, age)

	// Добавляем заголовки CORS для основного запроса
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := JsonResponse{
		Status:  "success",
		Message: fmt.Sprintf("Name: %s, Age: %s", name, age),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}

//func saveRegistrationData(username, email, password string) error {
//	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", username, email, password)
//	return err
//}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	os.Setenv("DATABASE_URL", "user=postgres dbname=postgres password=admin sslmode=disable")
	var err error
	db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User{}) // Создание таблицы User (если её нет)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")

	// db.Create(&User{Username: "Sultan", Email: "abdukarimov.05@gmail.com", Password: "asdsad&asd"})
	// db.Model(&User{}).Where("id = ?", 2).Update("username", "Aaaaaa")
	// db.Unscoped().Delete(&User{}, "id = ?", 3)

	// user, err := getUserByID(4)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("User by ID:", user)
	// }

	// allUsers, err := getAllUsers()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("All users:", allUsers)
	// }

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

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

func main() {
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
		fmt.Printf("Received POST request:\nUsername: %s\nEmail: %s\nPassword: %s\n --------------\n", requestData.Username, requestData.Email, requestData.Password)
	}

	saveError := saveRegistrationData(requestData.Username, requestData.Email, requestData.Password)
	if saveError != nil {
		http.Error(w, `{"status": "500", "message": "Failed to save registration data"}`, http.StatusInternalServerError)
		return
	}

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

func saveRegistrationData(username, email, password string) error {
	user := User{Username: username, Email: email, Password: password}
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func getUserByID(userID uint) (User, error) {
	var user User
	result := db.Select("id, username, email, password").First(&user, userID)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func getAllUsers() ([]User, error) {
	var users []User
	result := db.Select("id, username, email, password").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

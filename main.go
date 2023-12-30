package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonRequest struct {
	Message string `json:"message"`
}

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", handlePostRequest)
	http.ListenAndServe(":8080", nil)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData JsonRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, `{"status": "400", "message": "Некорректное JSON-сообщение"}`, http.StatusBadRequest)
		return
	}

	if requestData.Message == "" {
		http.Error(w, `{"status": "400", "message": "Некорректное JSON-сообщение"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("Received message:", requestData.Message)

	response := JsonResponse{
		Status:  "success",
		Message: "Данные успешно приняты",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

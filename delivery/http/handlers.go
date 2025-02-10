package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"gitverse-internship-zg/services/user-service/domain/entities"
	"gitverse-internship-zg/services/user-service/usecase"
)

func UsersHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetUser(uc, w, r)
		case "POST":
			CreateUser(uc, w, r)
		case "PATCH":
			EditUser(uc, w, r)
		default:
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	}
}

func GetUser(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	slog.Info("handle get user id=" + id)
	// id := "52d39105-96c2-408c-91c9-05ee1d4b99e7"
	entity, err := uc.GetEntityByID(id)
	fmt.Println(entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func EditUser(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {
	var newUserInfo entities.EntityRequest
	fmt.Println("handler", newUserInfo)
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&newUserInfo); err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("handler2", newUserInfo)
	defer r.Body.Close()

	answer, err := uc.EditUserByID(newUserInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}

func CreateUser(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {
	var newUserInfo entities.EntityRequest
	fmt.Println("handler", newUserInfo)
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&newUserInfo); err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("handler2", newUserInfo)
	defer r.Body.Close()

	answer, err := uc.CreateUser(newUserInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}

package ht

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"simple_api/internal/repository"
	"simple_api/internal/service"

	"github.com/gorilla/mux"
)

// UserHandler struct для работы с пользователями
type UserHandler struct {
	repo *repository.UserRepository
}

// NewHandler инициализирует маршруты
func NewHandler(r *mux.Router, db *sql.DB) {
	repo := repository.NewUserRepository(db)
	handler := &UserHandler{repo: repo}

	// Создаем маршруты
	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	// Защищенные маршруты
	r.HandleFunc("/users/{id:[0-9]+}", handler.GetUserByID).Methods("GET").Handler(AuthMiddleware(http.HandlerFunc(handler.GetUserByID)))
	r.HandleFunc("/users/{id:[0-9]+}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", handler.DeleteUser).Methods("DELETE")
    }




func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user service.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	
	if user.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}


	createdUser, err := h.repo.Create(user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерируем токен
	token, err := service.GenerateToken(createdUser.ID)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен и данные пользователя
	response := map[string]interface{}{
		"user":  createdUser,
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


// Получение всех пользователей
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}


func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Извлекаем user_id из контекста
	userID := r.Context().Value("user_id").(int)

	// Если пользователь пытается получить данные не о себе, возвращаем ошибку
	id := mux.Vars(r)["id"]
	if id != fmt.Sprint(userID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Дальше идет обычная логика
	user, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}


func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user service.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.repo.Update(id, user.Name,user.Password, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}


func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := h.repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
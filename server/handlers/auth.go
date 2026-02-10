package handlers

import (
	"log"
	"net/http"
	"status-page-monitor/internal/database"
	"status-page-monitor/internal/database/models"
	"status-page-monitor/server/jwt"
	res "status-page-monitor/server/response"
	"status-page-monitor/server/utils"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if ok := utils.DecodeJSON(w, r, &req); ok != true {
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	var user models.User

	result := database.DB.Where("login = ?", req.Login).First(&user)

	if result.Error != nil {
		log.Printf("User not found: %s", req.Login)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		log.Printf("Invalid password for user: %s", req.Login)
		log.Printf("DB password: %s", user.Password)
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		log.Printf("REQ password: %s", string(hash))
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateJWT(user)

	if err != nil {
		log.Println("Unabler to create JWT")
		http.Error(w, "Internal", http.StatusInternalServerError)
	}

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	response.User.Password = ""

	w.Header().Set("Content-Type", "application/json")
	res.JSON(w, http.StatusOK, response)
}

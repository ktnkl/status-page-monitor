package handlers

import (
	"log"
	"net/http"
	"status-page-monitor/internal/database"
	"status-page-monitor/internal/database/models"
	"status-page-monitor/server/response"
	"status-page-monitor/server/utils"
	"time"
)

type CreateServerRequest struct {
	Url      string `json:"url" validate:"required,url"`
	Interval int    `json:"interval" validate:"required"`
}

func GetAllServersHandler(w http.ResponseWriter, r *http.Request) {
	var servers []models.Server
	var total int64
	var pagination utils.Pagination

	result := database.DB.Scopes(utils.Paginate(r, &pagination)).Find(&servers).Count(&total)

	if result.Error != nil {
		log.Println("Error while fetching servers from db:", result.Error)
	}

	pagination.Total = total

	res := utils.PaginatedResponse{Pagination: pagination, Data: servers}

	response.OK(w, res)
	// pagination := utils.Pagination{Page: }
}

func GetServerById(w http.ResponseWriter, r *http.Request) {
	a := r.Context()
	issuer := a.Value("login")

	response.Success(w, http.StatusOK, issuer)
}

func CreateServer(w http.ResponseWriter, r *http.Request) {
	var req CreateServerRequest

	if ok := utils.DecodeJSON(w, r, &req); ok != true {
		return
	}

	if req.Url == "" || req.Interval == 0 {
		details := make(map[string]string)
		details["Required"] = "Url and interval are required"
		response.ValidationError(w, details)
	}

	server := models.Server{Url: req.Url, Interval: req.Interval, Checkedat: time.Now().UTC().Format(time.RFC3339), Nextcheckat: time.Now().Add(time.Duration(req.Interval) * time.Minute).UTC().Format(time.RFC3339)}

	result := database.DB.Create(&server)

	if result.Error != nil {
		log.Println("Error while creating server: ", result.Error)
		response.InternalError(w)
	}

	response.Created(w, server)
}

func EditServer(w http.ResponseWriter, r *http.Request) {
	a := r.Context()
	issuer := a.Value("login")

	response.Success(w, http.StatusOK, issuer)
}

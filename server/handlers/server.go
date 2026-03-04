package handlers

import (
	"log"
	"net/http"
	"status-page-monitor/internal/database"
	"status-page-monitor/internal/database/models"
	"status-page-monitor/server/response"
	"status-page-monitor/server/utils"
	"time"

	"github.com/gorilla/mux"
)

type CreateServerRequest struct {
	Url      string `json:"url" validate:"required,url"`
	Interval int    `json:"interval" validate:"required"`
}

func GetAllServersHandler(w http.ResponseWriter, r *http.Request) {
	var servers []models.Server
	var total int64
	var pagination utils.Pagination

	result := database.DB.Scopes(utils.Paginate(r, &pagination)).Find(&servers)

	if result.Error != nil {
		log.Println("Error while fetching servers from db:", result.Error)
	}

	if err := database.DB.Table("servers").Count(&total).Error; err != nil {
		log.Println("Error while getting total count in GET_SERVERS_ALL:", err)
	}

	pagination.Total = total

	res := utils.PaginatedResponse{Pagination: pagination, Data: servers}

	response.OK(w, res)
}

func GetServerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var server models.Server

	if vars["id"] == "" {
		log.Println("Invalid id in get server by id:")
		response.InternalError(w)
		return
	}

	if err := database.DB.First(&server, vars["id"]).Error; err != nil {
		log.Println("Error while finding server by id:", err)
		response.IDNotFound(w, vars["id"])
		return
	}

	response.OK(w, server)
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
		return
	}

	server := models.Server{Url: req.Url, Interval: req.Interval, Checkedat: time.Now().UTC().Format(time.RFC3339), Nextcheckat: time.Now().Add(time.Duration(req.Interval) * time.Second).UTC().Format(time.RFC3339)}

	result := database.DB.Create(&server)

	if result.Error != nil {
		log.Println("Error while creating server: ", result.Error)
		response.InternalError(w)
		return
	}

	response.Created(w, server)
}

func EditServer(w http.ResponseWriter, r *http.Request) {
	a := r.Context()
	issuer := a.Value("login")

	response.Success(w, http.StatusOK, issuer)
}

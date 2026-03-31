package todo

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MayurUbarhande0/backend/internals/middleware"
	"github.com/gorilla/mux"
)

var taskRepo *Mydb

type createTaskRequest struct {
	Title string `json:"title"`
}

type updateTaskRequest struct {
	Done bool `json:"done"`
}

func SetDB(db *sql.DB) {
	taskRepo = NewMydb(db)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("title is required")
		return
	}

	task, err := taskRepo.Addtask(userID, req.Title)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")
		return
	}

	id := mux.Vars(r)["id"]
	taskID, err := strconv.Atoi(id)
	if err != nil || taskID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid task id")
		return
	}

	err = taskRepo.Deletetask(taskID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("task deleted")
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")
		return
	}

	id := mux.Vars(r)["id"]
	taskID, err := strconv.Atoi(id)
	if err != nil || taskID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid task id")
		return
	}

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}

	task, err := taskRepo.UpdateTask(taskID, userID, req.Done)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")
		return
	}

	tasks, err := taskRepo.GetAll(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

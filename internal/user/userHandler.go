package user

import (
	"encoding/json"
	"go-tasklist/internal/task"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Api interface {
	GetUserTasks(res http.ResponseWriter, req *http.Request)
	AddUserTask(res http.ResponseWriter, req *http.Request)
	DeleteUserTask(res http.ResponseWriter, req *http.Request)
}

type api struct {
	s UserService
}

type TaskRequest struct {
	UserId int64
	TaskId *int64
}

func parseRequestPathValues(r *http.Request) (TaskRequest, error) {
	userIdRaw := r.PathValue("userId")
	taskIdRaw := r.PathValue("taskId")
	userId, err := strconv.ParseInt(userIdRaw, 10, 64)
	if nil != err {
		return TaskRequest{}, err
	}
	if (len(taskIdRaw) > 0) {
		taskId, err := strconv.ParseInt(taskIdRaw, 10, 64)
		if nil != err {
			return TaskRequest{}, err
		}
		return TaskRequest{
			UserId: userId,
			TaskId: &taskId,
		}, nil
	}
	return TaskRequest{
		UserId: userId,
		TaskId: nil,
	}, nil
}

func parseRequestBody(r *http.Request) (task.TaskDto, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if nil != err {
		return task.TaskDto{}, err
	}

	var newTask task.TaskDto

	err = json.Unmarshal(bodyBytes, &newTask)
	if nil != err {
		return task.TaskDto{}, err
	}

	return newTask, nil

}

// route GET /api/tasks/{userId}
func (a api) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	reqParams, err := parseRequestPathValues(r)
	if (nil != err) {
		log.Fatalf("What is this thing?  %s\n", err)
	}
	tasks := a.s.GetTasks(reqParams.UserId)
	response := make(map[string]any)
	response["status"] = "success"
	response["data"] = tasks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// route POST /api/tasks/{userId}
func (a api) AddUserTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	reqParams, err := parseRequestPathValues(r)
	if (nil != err) {
		log.Fatalf("What is this thing?  %s\n", err)
	}

	newTask, err := parseRequestBody(r)
	if (nil != err) {
		log.Fatalf("What is this thing?  %s\n", err)
	}

	a.s.AddTask(reqParams.UserId, newTask)

	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	response["status"] = "success"
	if nil != err {
		log.Printf("AddUserTask Marshal")
		log.Fatal(err)
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)
}

// route DELETE /api/tasks/{userId}
func (a api) DeleteUserTask(w http.ResponseWriter, r *http.Request) {
	reqParams, err := parseRequestPathValues(r)
	if (nil != err) {
		log.Fatalf("What is this thing?  %s\n", err)
	}

	a.s.RemoveTask(reqParams.UserId, *reqParams.TaskId)

	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	response["status"] = "success"
	json.NewEncoder(w).Encode(response)
}

func GetUserAPI(s UserService) Api {
	return api{s}
}

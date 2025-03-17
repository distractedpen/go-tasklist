package task

import (
	"encoding/json"
	"go-tasklist/internal/util"
	"io"
	"log"
	"net/http"
	"strconv"
)
type Api interface {
	GetUserTasks(http.ResponseWriter, *http.Request)
	AddUserTask(http.ResponseWriter, *http.Request)
	DeleteUserTask(http.ResponseWriter, *http.Request)
}

type api struct {
	s TaskService
}

type TaskRequest struct {
	UserId int64
	TaskId *int64 // task id may not alway be given
}

func parseRequestPathValues(r *http.Request) (TaskRequest, error) {
	userIdRaw := r.PathValue("userId")
	taskIdRaw := r.PathValue("taskId")
	userId, err := strconv.ParseInt(userIdRaw, 10, 64)
	if nil != err {
		return TaskRequest{}, util.ErrRequestInvalid{}
	}

	if len(taskIdRaw) > 0 {
		taskId, err := strconv.ParseInt(taskIdRaw, 10, 64)
		if nil != err {
			return TaskRequest{}, util.ErrRequestInvalid{}
		}
		return TaskRequest{
			UserId: userId,
			TaskId: &taskId,
		}, nil
	} else {
		return TaskRequest{
			UserId: userId,
			TaskId: nil,
		}, nil
	}
}

func parseRequestBody(r *http.Request) (TaskDto, error) {
	var newTask TaskDto
	bodyBytes, err := io.ReadAll(r.Body)
	if nil != err {
		return newTask, util.ErrRequestInvalid{}
	}

	err = json.Unmarshal(bodyBytes, &newTask)
	if nil != err {
		return newTask, util.ErrRequestInvalid{}
	}

	return newTask, nil
}

// route GET /api/tasks/{userId}
func (a api) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	reqParams, err := parseRequestPathValues(r)
	if nil != err {
		log.Printf("What is this thing?  %v\n", err)
		util.SendResponse(w, 400, nil)
		return
	}

	tasks, err := a.s.GetTasks(reqParams.UserId)
	if nil != err {
		log.Printf("I broke it!  %v\n", err)
		util.SendResponse(w, 500, nil)
		return
	}
	util.SendResponse(w, 200, tasks)
}

// route POST /api/tasks/{userId}
func (a api) AddUserTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	reqParams, err := parseRequestPathValues(r)
	if nil != err {
		log.Printf("What is this thing?  %v\n", err)
		util.SendResponse(w, 400, nil)
		return
	}

	newTask, err := parseRequestBody(r)
	if nil != err {
		log.Printf("What is this thing?  %v\n", err)
		util.SendResponse(w, 400, nil)
		return
	}

	err = a.s.AddTask(reqParams.UserId, newTask)
	if nil != err {
		log.Printf("I broke it!  %v\n", err)
		util.SendResponse(w, 500, nil)
		return
	}

	util.SendResponse(w, 201, nil)
}

// route DELETE /api/tasks/{userId}
func (a api) DeleteUserTask(w http.ResponseWriter, r *http.Request) {
	reqParams, err := parseRequestPathValues(r)
	if nil != err {
		log.Printf("What is this thing?  %v\n", err)
		util.SendResponse(w, 400, nil)
		return
	}

	err = a.s.RemoveTask(reqParams.UserId, *reqParams.TaskId)
	if nil != err {
		log.Printf("I broke it!  %v\n", err)
		util.SendResponse(w, 500, nil)
		return
	}

	util.SendResponse(w, 204, nil)
}

func GetTaskApi(s TaskService) Api {
	return api{s}
}

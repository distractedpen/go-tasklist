package user

import (
	"log"
	"go-tasklist/internal/task"
)

// Could use a custom error type to handle Does not exist or other errors

// User Service
type UserService interface {
	GetTasks(int64) []task.Task
	AddTask(int64, task.TaskDto)
	RemoveTask(int64, int64)
}

type userService struct {
	taskRepository task.TaskRepository
}

// Get Task
func (s userService) GetTasks(userId int64) []task.Task {
	tasks, err := s.taskRepository.GetTasks(userId)
	if nil != err {
		log.Fatal(err)
	}
	return tasks
}

// Add Task
func (s userService) AddTask(userId int64, newTask task.TaskDto) {
	err := s.taskRepository.InsertTask(userId, newTask)
	if nil != err {
		log.Fatal(err)
	}
}

// Remove Task
func (s userService) RemoveTask(userId int64, taskId int64) {
	err := s.taskRepository.RemoveTask(userId, taskId)
	if nil != err {
		log.Fatal(err)
	}
}

func GetUserService(taskRepository task.TaskRepository) UserService {
	return userService{
		taskRepository: taskRepository,
	} }

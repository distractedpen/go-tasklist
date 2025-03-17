package task

type TaskService interface {
	GetTasks(int64) ([]Task, error)
	AddTask(int64, TaskDto) error
	RemoveTask(int64, int64) error
}

type taskService struct {
	taskRepository TaskRepository
}

func (s taskService) GetTasks(userId int64) ([]Task, error) {
	tasks, err := s.taskRepository.GetTasks(userId)
	return tasks, err
}

func (s taskService) AddTask(userId int64, newTask TaskDto) error {
	err := s.taskRepository.InsertTask(userId, newTask)
	return err
}

func (s taskService) RemoveTask(userId int64, taskId int64) error {
	err := s.taskRepository.RemoveTask(userId, taskId)
	return err
}

func GetTaskService(taskRepository TaskRepository) TaskService {
	return taskService{
		taskRepository: taskRepository,
	}
}

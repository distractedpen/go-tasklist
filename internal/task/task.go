package task

type Task struct {
	Id int64           `json:"id"`
	Name string        `json:"name"`
	Description string `json:"description"`
}

type TaskDto struct {
	Name string        `json:"name"`
	Description string `json:"description"`
}

func (t *TaskDto) MakeTask(id int64) Task {
	return Task{ Id: id,
		Name: t.Name,
		Description: t.Description,
	}
}


func (t *Task) MakeDto() TaskDto {
	return TaskDto{
		Name: t.Name,
		Description: t.Description,
	}
}

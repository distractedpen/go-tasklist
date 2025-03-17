package task

import (
	"database/sql"
	"log"

	"go-tasklist/internal/util"
)

type TaskRepository interface {
	GetTasks(int64) ([]Task, error)
	InsertTask(int64, TaskDto) error
	RemoveTask(int64, int64) error
}

type taskRepository struct {
	db *sql.DB
}

func (r taskRepository) GetTasks(userId int64) ([]Task, error) {
	tasks := make([]Task, 0)

	rows, err := r.db.Query(`
		SELECT * FROM tasks where id IN 
		(SELECT taskId FROM users_tasks WHERE userId = ?);`, 
		userId)
	defer rows.Close()

	if nil != err {
		return tasks, util.ErrDB{}
	}

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Name, &task.Description)
		if nil != err {
			return tasks, util.ErrDB{}
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r taskRepository) InsertTask(userId int64, newTask TaskDto) error {
	row := r.db.QueryRow(`
		SELECT * FROM tasks WHERE name=? AND description=?;`,
		newTask.Name, newTask.Description,
	)

	var task Task
	err := row.Scan(&task.Id, &task.Name, &task.Description)

	if nil != err && sql.ErrNoRows != err {
		log.Println(err.Error())
		return util.ErrDB{}
	}
	newId := task.Id

	if sql.ErrNoRows == err {
		// Item does not exist in task table

		// check if user list already contains this task
		row := r.db.QueryRow(
			`SELECT * FROM users_tasks WHERE userId=? AND taskId=?`,
			userId, task.Id,
		)

		if row.Scan() != sql.ErrNoRows {
			return util.ErrExists{}
		}

		res, err := r.db.Exec(
			"INSERT INTO tasks(name, description) values (?, ?);",
			newTask.Name, newTask.Description)

		if nil != err {
			log.Println(err.Error())
			return util.ErrDB{}
		}

		newId, _ = res.LastInsertId()
	}

	_, err = r.db.Exec("INSERT INTO users_tasks values (?,?);", userId, newId)
	if nil != err {
		log.Println(err.Error())
		return util.ErrDB{}
	}

	return nil
}

// Further Idea: once no user references a task, remove it from task table
func (r taskRepository) RemoveTask(userId int64, taskId int64) error {
	row := r.db.QueryRow(
		`SELECT * FROM users_tasks WHERE userId=? AND taskId=?;`,
		userId, taskId)

	err := row.Err()
	if nil != err && sql.ErrNoRows != err {
		log.Println(err.Error())
		return util.ErrDB{}
	}

	if sql.ErrNoRows == err {
		return util.ErrDoesNotExist{}
	}

	_, err = r.db.Exec(
		`DELETE FROM users_tasks WHERE userId=? AND taskId=?;`,
		userId, taskId)
	if nil != err {
		log.Print(err)
		return util.ErrDB{}
	}

	return nil
}

func GetTaskRepository(db *sql.DB) TaskRepository {
	return taskRepository{db}
}

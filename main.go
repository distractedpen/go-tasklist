package main

import (
	"log"
	"net/http"
	"time"

	"go-tasklist/internal/task"
	"go-tasklist/internal/user"

	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func getDB() *sql.DB {
	var version string
	db, err := sql.Open("sqlite3", "file:tasklist.db")
	if nil != err {
		log.Fatal(err)
	}
	db.QueryRow(`SELECT sqlite_version()`).Scan(&version)
	log.Printf("Sqlite3 %s\n", version)
	return db
}

func createTables(db *sql.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS
		tasks
		( id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name TEXT, 
			description TEXT
		);`)

	db.Exec(`CREATE TABLE IF NOT EXISTS
	    users
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			username TEXT,
			password TEXT
		);`)

	db.Exec(`CREATE TABLE IF NOT EXISTS
	    users_tasks
		(
			userId INTEGER, 
			taskId INTEGER
		);`)
}

type uiData struct {
	User user.UserDto
	Tasks []task.Task
}

func main() {

	db := getDB()
	createTables(db)

	mux := http.NewServeMux()

	staticFS := http.FileServer(http.Dir("./web/static/"))
	taskRespoitory := task.GetTaskRepository(db)
	userService := user.GetUserService(taskRespoitory)
	userApi := user.GetUserAPI(userService)

	log.Printf("Loaded userApi: %+v", userApi)

	mux.Handle("/", staticFS)

	mux.HandleFunc("GET /api/tasks/{userId}", userApi.GetUserTasks)
	mux.HandleFunc("POST /api/tasks/{userId}", userApi.AddUserTask)
	mux.HandleFunc("DELETE /api/tasks/{userId}/{taskId}", userApi.DeleteUserTask)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening on :8080")
	log.Fatal(s.ListenAndServe())
}

package main

import (
	"log"
	"net/http"
	"time"

	"go-tasklist/internal/auth"
	"go-tasklist/internal/task"
	"go-tasklist/internal/user"
	"go-tasklist/internal/middleware"

	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/joho/godotenv"
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

func main() {
	err := godotenv.Load(".env")
	if nil != err {
		log.Fatalf(".env file not found\n")
	}


	db := getDB()
	createTables(db)

	mux := http.NewServeMux()

	staticFS := http.FileServer(http.Dir("./web/static/"))
	taskRespoitory := task.GetTaskRepository(db)
	tasksService := task.GetTaskService(taskRespoitory)
	tasksApi := task.GetTaskApi(tasksService)
	userRepository := user.GetUserRepository(db)
	userService := user.GetUserService(userRepository)
	authService := auth.GetAuthService(userService)
	authApi := auth.GetAuthHandlers(authService)

	log.Printf("Loaded tasksApi: %+v", tasksApi)
	log.Printf("Loaded staticFS: %+v", staticFS)

	mux.Handle("/", staticFS);

	mux.HandleFunc("POST /api/auth/login", authApi.Login)
	mux.HandleFunc("POST /api/auth/register", authApi.Register)

	mux.Handle("GET /api/tasks/{userId}", 
		middleware.Authenticated(http.HandlerFunc(tasksApi.GetUserTasks)))
	mux.Handle("POST /api/tasks/{userId}", 
		middleware.Authenticated(http.HandlerFunc(tasksApi.AddUserTask)))
	mux.Handle("DELETE /api/tasks/{userId}/{taskId}", 
		middleware.Authenticated(http.HandlerFunc(tasksApi.DeleteUserTask)))

	loggerMux := middleware.NewLoggerHandler(mux)


	s := &http.Server{
		Addr:           ":8080",
		Handler:        loggerMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening on :8080")
	log.Fatal(s.ListenAndServe())
}

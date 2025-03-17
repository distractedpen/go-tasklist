package middleware

import (
	"log"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
}


func NewLoggerHandler(h http.Handler) *Logger {
	return &Logger{h}
}

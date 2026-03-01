package api

import "net/http"

const webDir = "./web"

func Init() {
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)

	http.Handle("/", http.FileServer(http.Dir(webDir)))
}

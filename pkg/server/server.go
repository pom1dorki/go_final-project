package server

import (
	"GO_FINAL_PROJECT/pkg/api"
	"fmt"
	"net/http"
	"os"
)

func New() {
	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	addr := ":" + port

	webDir := "web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

package main

import (
	"GO_FINAL_PROJECT/pkg/db"
	"GO_FINAL_PROJECT/pkg/server"
	"fmt"
	"os"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		fmt.Println("Ошибка инициализации базы данных:", err)
		return
	}
	fmt.Println("База данных успешно инициализирована (scheduler.db)")

	server.New()
}

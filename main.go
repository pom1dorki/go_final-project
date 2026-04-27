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
	fmt.Println("База данных успешно инициализирована:", dbFile)

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Ошибка закрытия базы данных:", err)
		}
	}()

	server.New()
}

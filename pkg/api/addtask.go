package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"GO_FINAL_PROJECT/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour)

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты")
	}
	t = t.Truncate(24 * time.Hour)

	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		if t.Before(now) {
			task.Date, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if t.Before(now) {
		task.Date = now.Format(DateFormat)
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

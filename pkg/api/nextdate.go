package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("cannot parse dstart: %w", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}

	rule := parts[0]

	switch rule {
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("d rule requires number")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", fmt.Errorf("invalid d interval")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("unsupported repeat rule")
	}
}

func afterNow(date, now time.Time) bool {
	d := date.Truncate(24 * time.Hour)
	n := now.Truncate(24 * time.Hour)
	return d.After(n)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "cannot parse now", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}

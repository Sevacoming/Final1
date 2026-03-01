package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"final1/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var task db.Task

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	task.Date = strings.TrimSpace(task.Date)
	task.Title = strings.TrimSpace(task.Title)
	task.Comment = strings.TrimSpace(task.Comment)
	task.Repeat = strings.TrimSpace(task.Repeat)

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	today := dateOnly(time.Now())
	todayStr := today.Format("20060102")

	if task.Repeat != "" {
		base := task.Date
		if base == "" {
			base = todayStr
		}
		if _, err := NextDate(today, base, task.Repeat); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}

	if task.Date == "" {
		task.Date = todayStr
	} else {
		d, err := time.ParseInLocation("20060102", task.Date, time.Local)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
			return
		}
		d = dateOnly(d)

		if d.Before(today) {
			if task.Repeat == "" {

				task.Date = todayStr
			} else {

				next, err := NextDate(today, task.Date, task.Repeat)
				if err != nil {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				task.Date = next
			}
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}

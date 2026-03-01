package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"final1/pkg/db"
)

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	t.ID = strings.TrimSpace(t.ID)
	t.Date = strings.TrimSpace(t.Date)
	t.Title = strings.TrimSpace(t.Title)
	t.Comment = strings.TrimSpace(t.Comment)
	t.Repeat = strings.TrimSpace(t.Repeat)

	if t.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if t.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок"})
		return
	}

	if _, err := time.Parse("20060102", t.Date); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат даты"})
		return
	}

	if t.Repeat != "" {
		if _, err := NextDate(time.Now(), t.Date, t.Repeat); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат повтора"})
			return
		}
	}

	if err := db.UpdateTask(&t); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{}) // {}
}

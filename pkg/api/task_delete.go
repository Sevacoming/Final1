package api

import (
 "net/http"
 "strings"

 "final1/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimSpace(r.URL.Query().Get("id"))
 if id == "" {
  writeJSON(w, http.StatusOK, map[string]string{"error": "Не указан идентификатор"})
  return
 }

 if err := db.DeleteTask(id); err != nil {
  writeJSON(w, http.StatusOK, map[string]string{"error": err.Error()})
  return
 }

 writeJSON(w, http.StatusOK, map[string]any{})
}

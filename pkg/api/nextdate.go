package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	start, err := time.ParseInLocation("20060102", strings.TrimSpace(dstart), time.Local)
	if err != nil {
		return "", err
	}

	now = dateOnly(now)
	start = dateOnly(start)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid repeat format")
		}
		start = start.AddDate(1, 0, 0)
		for !start.After(now) {
			start = start.AddDate(1, 0, 0)
		}
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid repeat format")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", errors.New("invalid day interval")
		}
		start = start.AddDate(0, 0, n)
		for !start.After(now) {
			start = start.AddDate(0, 0, n)
		}
	default:
		return "", errors.New("unsupported repeat format")
	}

	return start.Format("20060102"), nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := strings.TrimSpace(r.FormValue("now"))
	dateStr := strings.TrimSpace(r.FormValue("date"))
	repeat := r.FormValue("repeat")

	now := time.Now()
	if nowStr != "" {
		t, err := time.ParseInLocation("20060102", nowStr, time.Local)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = t
	}

	if dateStr == "" {
		http.Error(w, "date is empty", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	_, _ = w.Write([]byte(next))
}

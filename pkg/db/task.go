package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, errors.New("db not initialized")
	}
	res, err := DB.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 50 {
		limit = 50
	}

	today := time.Now().Format("20060102")

	// базовый запрос
	query := `SELECT id, date, title, comment, repeat
FROM scheduler
WHERE date >= ?
ORDER BY date ASC
LIMIT ?`

	args := []any{today, limit}

	if search != "" {
		query = `SELECT id, date, title, comment, repeat
FROM scheduler
WHERE date >= ?
  AND (title LIKE ? OR comment LIKE ?)
ORDER BY date ASC
LIMIT ?`
		p := "%" + search + "%"
		args = []any{today, p, p, limit}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*Task, 0, limit)

	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}

	t := &Task{}
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return t, nil
}

func UpdateTask(task *Task) error {
	if DB == nil {
		return errors.New("db not initialized")
	}

	res, err := DB.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func DeleteTask(id string) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

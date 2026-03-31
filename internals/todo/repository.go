package todo

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewMydb(db *sql.DB) *Mydb {
	return &Mydb{db: db}
}

func (db *Mydb) Addtask(user_id int, title string) (Task, error) {
	result, err := db.db.Exec("INSERT INTO tasks (user_id, title) VALUES (?, ?)", user_id, title)
	if err != nil {
		return Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	var task Task
	err = db.db.QueryRow("SELECT id, user_id, title, done, created_at FROM tasks WHERE id = ?", id).
		Scan(&task.ID, &task.User_id, &task.Title, &task.Done, &task.CreatedAt)
	if err != nil {
		return Task{}, err
	}

	return task, err
}

func (db *Mydb) GetAll(userID int) ([]Task, error) {
	rows, err := db.db.Query("SELECT id, user_id, title, done, created_at FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.User_id, &task.Title, &task.Done, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *Mydb) Deletetask(id int, userID int) error {
	_, err := db.db.Exec("DELETE FROM tasks WHERE id = ? AND user_id = ?", id, userID)
	return err
}

func (db *Mydb) UpdateTask(id int, userID int, what bool) (Task, error) {
	_, err := db.db.Exec("UPDATE tasks SET done = ? WHERE id = ? AND user_id = ?", what, id, userID)
	if err != nil {
		return Task{}, err
	}

	var task Task
	err = db.db.QueryRow("SELECT id, user_id, title, done, created_at FROM tasks WHERE id = ?", id).
		Scan(&task.ID, &task.User_id, &task.Title, &task.Done, &task.CreatedAt)

	return task, err
}

package postsql

import (
	"database/sql"
	"fmt"
	"golangify.com/SnippetBox/pkg/models"
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type SnippetModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// Ниже будет SQL запрос, который мы хотим выполнить. Мы разделили его на две строки
	// для удобства чтения (поэтому он окружен обратными кавычками
	// вместо обычных двойных кавычек).
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, current_timestamp, current_timestamp + interval '1 year'))`
	_, err := m.DB.Exec(stmt, title, content, expires) // много вопросов
	if err != nil {
		return 0, err
	}
	// Используем метод LastInsertId(), чтобы получить последний ID
	// созданной записи из таблицу snippets.
	var id int
	err = m.DB.QueryRow("INSERT INTO snippets (title) VALUES ('John') RETURNING id").Scan(&id)
	if err != nil {
		return 0, err

	}
	fmt.Println(id)
	return int(id), err // id - int64
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}

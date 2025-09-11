package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // драйвер PostgreSQL
)

var DB *sql.DB

// Инициализация подключения
func InitDB() error {
	connStr := "host=localhost port=5432 user=postgres password=120311 dbname=forumdb sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка открытия подключения: %w", err)
	}

	// проверим подключение
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// минимальная миграция: колонка для хранения изображения в постах
	if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS image_data BYTEA`); err != nil {
		return fmt.Errorf("ошибка миграции (image_data): %w", err)
	}

	// добавить недостающие колонки под текущую модель
	if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS image_url TEXT`); err != nil {
		return fmt.Errorf("ошибка миграции (image_url): %w", err)
	}
	if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS link_url TEXT`); err != nil {
		return fmt.Errorf("ошибка миграции (link_url): %w", err)
	}
	if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS board_id INTEGER`); err != nil {
		return fmt.Errorf("ошибка миграции (board_id): %w", err)
	}
	if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now()`); err != nil {
		return fmt.Errorf("ошибка миграции (created_at): %w", err)
	}
	if _, err := DB.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`); err != nil {
		return fmt.Errorf("ошибка миграции (updated_at): %w", err)
	}

	// комментарии
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("ошибка миграции (comments): %w", err)
	}

	// лайки/дизлайки постов
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS post_votes (
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			value SMALLINT NOT NULL CHECK (value IN (-1, 1)),
			PRIMARY KEY (post_id, user_id)
		)`); err != nil {
		return fmt.Errorf("ошибка миграции (post_votes): %w", err)
	}

	// лайки комментариев
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS comment_votes (
			comment_id INTEGER NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			value SMALLINT NOT NULL CHECK (value IN (-1, 1)),
			PRIMARY KEY (comment_id, user_id)
		)`); err != nil {
		return fmt.Errorf("ошибка миграции (comment_votes): %w", err)
	}

	// уникальные просмотры постов
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS post_views (
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			PRIMARY KEY (post_id, user_id)
		)`); err != nil {
		return fmt.Errorf("ошибка миграции (post_views): %w", err)
	}

	fmt.Println("Подключение к базе успешно")
	return nil
}

// Закрытие подключения
func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Подключение к базе закрыто")
	}
}

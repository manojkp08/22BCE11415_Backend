package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connectionString string) error {
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL database")
	return nil
}

func GetOrCreateUser(email, name string) (*User, error) {
	var user User
	err := DB.QueryRow("SELECT id, email, name, created_at FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Email, &user.Name, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new user
		user.ID = generateUUID()
		user.Email = email
		user.Name = name
		user.CreatedAt = time.Now()

		_, err := DB.Exec(
			"INSERT INTO users (id, email, name, created_at) VALUES ($1, $2, $3, $4)",
			user.ID, user.Email, user.Name, user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		return &user, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func generateUUID() string {
	return uuid.New().String()
}

func CreateFile(file File) (*File, error) {
	_, err := DB.Exec(
		"INSERT INTO files (id, user_id, name, path, size, mime_type, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		file.ID, file.UserID, file.Name, file.Path, file.Size, file.MimeType, time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func GetFilesByUserID(userID string) ([]File, error) {
	rows, err := DB.Query(
		"SELECT id, user_id, name, path, size, mime_type, created_at, is_public FROM files WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var file File
		err := rows.Scan(
			&file.ID, &file.UserID, &file.Name, &file.Path, &file.Size,
			&file.MimeType, &file.CreatedAt, &file.IsPublic,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func GetFileByID(fileID string) (*File, error) {
	var file File
	err := DB.QueryRow(
		"SELECT id, user_id, name, path, size, mime_type, created_at, is_public FROM files WHERE id = $1",
		fileID,
	).Scan(
		&file.ID, &file.UserID, &file.Name, &file.Path, &file.Size,
		&file.MimeType, &file.CreatedAt, &file.IsPublic,
	)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

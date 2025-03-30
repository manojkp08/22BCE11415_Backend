package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/manojkp08/22BCE11415_Backend/internal/cache"
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

func GetUserByID(userID string) (*User, error) {
	var user User
	err := DB.QueryRow(
		"SELECT id, email, name, created_at FROM users WHERE id = $1",
		userID,
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
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

// GetExpiredFiles returns all files that are older than 7 days and not public
func GetExpiredFiles() ([]File, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, name, path, size, mime_type, created_at, is_public 
		FROM files 
		WHERE created_at < NOW() - INTERVAL '7 DAYS' 
		AND is_public = false`,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return []File{}, nil // Return empty slice if no expired files
		}
		log.Printf("Error querying expired files: %v", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var expiredFiles []File
	for rows.Next() {
		var file File
		err := rows.Scan(
			&file.ID,
			&file.UserID,
			&file.Name,
			&file.Path,
			&file.Size,
			&file.MimeType,
			&file.CreatedAt,
			&file.IsPublic,
		)
		if err != nil {
			log.Printf("Error scanning file row: %v", err)
			continue // Skip problematic rows but continue processing others
		}
		expiredFiles = append(expiredFiles, file)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return nil, err
	}

	return expiredFiles, nil
}

// DeleteFile removes a file record from the database
func DeleteFile(fileID string) error {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}

	// Try to delete the file
	result, err := tx.Exec("DELETE FROM files WHERE id = $1", fileID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting file %s: %v", fileID, err)
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Printf("Error checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return sql.ErrNoRows
	}

	// Invalidate cache if exists
	if err := cache.InvalidateFileCache(fileID); err != nil {
		log.Printf("Error invalidating cache for file %s: %v", fileID, err)
		// Don't fail the operation just because cache invalidation failed
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

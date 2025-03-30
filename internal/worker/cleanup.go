package worker

import (
	"log"
	"os"
	"time"

	"github.com/manojkp08/22BCE11415_Backend/internal/database"
)

func StartCleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running cleanup worker...")
		files, err := database.GetExpiredFiles()
		if err != nil {
			log.Printf("Error getting expired files: %v", err)
			continue
		}

		for _, file := range files {
			// Delete file from storage
			if err := os.Remove(file.Path); err != nil {
				log.Printf("Error deleting file %s: %v", file.ID, err)
				continue
			}

			// Delete from database
			if err := database.DeleteFile(file.ID); err != nil {
				log.Printf("Error deleting file metadata %s: %v", file.ID, err)
			}
		}
	}
}

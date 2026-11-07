package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Enable CORS for web interface
func EnableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, *")

	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

// Utility functions
func GenerateClientID() string {
	return fmt.Sprintf("client-%d", time.Now().UnixNano())
}

func GenerateMessageID() uint32 {
	return uint32(time.Now().UnixNano() & 0xFFFFFFFF)
}

func SaveExfilData(clientID string, data []byte) {
	filename := fmt.Sprintf("exfil-%s-%d.bin", clientID, time.Now().Unix())
	// Save to file implementation
	fmt.Printf("Saved exfil data to %s (%d bytes)\n", filename, len(data))
}

func SaveUploadedFile(clientID string, data []byte) error {
	uploadDir := "uploads"

	if err := os.MkdirAll(uploadDir, 0700); err != nil {
		return err
	}

	// TODO: filename should be the same as what it was on the client machine
	filename := filepath.Join(
		uploadDir,
		fmt.Sprintf("upload-%s-%d.bin", clientID, time.Now().UnixNano()),
	)

	// Save to file implementation

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Saved uploaded file to %s (%d bytes)\n", filename, len(data))
	return nil

}

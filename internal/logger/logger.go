package logger

import (
	"log"
	"os"
	"path/filepath"
)

var (
	logger *log.Logger
)

func init() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Open log file
	f, err := os.OpenFile(filepath.Join("logs", "godo.log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Create multi-writer to log to both file and stdout
	logger = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...interface{}) {
	logger.Printf("[INFO] "+format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Printf("[ERROR] "+format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Printf("[DEBUG] "+format, v...)
}

func Fatal(format string, v ...interface{}) {
	logger.Fatalf("[FATAL] "+format, v...)
}

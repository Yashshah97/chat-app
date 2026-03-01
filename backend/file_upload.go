package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	maxUploadSize = 50 * 1024 * 1024 // 50MB
	uploadDir     = "./uploads"
)

// FileInfo holds information about uploaded files
type FileInfo struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	URL       string    `json:"url"`
	UploadedAt time.Time `json:"uploaded_at"`
	UserID    string    `json:"user_id"`
}

// initializeUploadDirectory creates the uploads directory if it doesn't exist
func initializeUploadDirectory() error {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		log.Printf("Creating uploads directory: %s", uploadDir)
		return os.MkdirAll(uploadDir, 0755)
	}
	return nil
}

// uploadFileHandler handles file uploads
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("File upload request from %s", r.RemoteAddr)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with max size
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	if fileHeader.Size > maxUploadSize {
		log.Printf("File too large: %d bytes", fileHeader.Size)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	fileID := uuid.New().String()
	ext := filepath.Ext(fileHeader.Filename)
	filename := fileID + ext

	// Create file path
	filePath := filepath.Join(uploadDir, filename)

	// Create file on disk
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Error copying file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Prepare response
	fileInfo := FileInfo{
		ID:        fileID,
		Filename:  fileHeader.Filename,
		Size:      fileHeader.Size,
		MimeType:  fileHeader.Header.Get("Content-Type"),
		URL:       fmt.Sprintf("/api/files/download/%s", filename),
		UploadedAt: time.Now(),
		UserID:    r.Header.Get("X-User-ID"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"status":"success","data":%+v}`, fileInfo)
	log.Printf("File uploaded successfully: %s", filename)
}

// downloadFileHandler handles file downloads
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filename := params["filename"]

	log.Printf("File download request for: %s", filename)

	// Validate filename to prevent path traversal
	if filepath.Clean(filename) != filename {
		log.Printf("Invalid filename: %s", filename)
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File not found: %s", filePath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	http.ServeFile(w, r, filePath)
	log.Printf("File downloaded successfully: %s", filename)
}

// deleteFileHandler handles file deletion
func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filename := params["filename"]

	log.Printf("File deletion request for: %s", filename)

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate filename
	if filepath.Clean(filename) != filename {
		log.Printf("Invalid filename: %s", filename)
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(uploadDir, filename)

	// Delete file
	if err := os.Remove(filePath); err != nil {
		log.Printf("Error deleting file: %v", err)
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting file", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"success","message":"File deleted"}`)
	log.Printf("File deleted successfully: %s", filename)
}

// listFilesHandler lists all uploaded files
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("List files request")

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Printf("Error reading upload directory: %v", err)
		http.Error(w, "Error reading files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "[")
	for i, file := range files {
		fileInfo, _ := file.Info()
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprintf(w, `{"name":"%s","size":%d,"modified":"%s"}`,
			file.Name(), fileInfo.Size(), fileInfo.ModTime().Format(time.RFC3339))
	}
	fmt.Fprint(w, "]")
}

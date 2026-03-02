package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// BatchOperation manages bulk operations
type BatchOperation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `json:"name"`
	Type          string         `json:"type"` // "delete", "archive", "export"
	Status        string         `json:"status"` // "pending", "in_progress", "completed"
	TotalItems    int            `json:"total_items"`
	ProcessedItems int           `json:"processed_items"`
	FailedItems   int            `json:"failed_items"`
	Parameters    datatypes.JSON `gorm:"type:jsonb" json:"parameters"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CompletedAt   time.Time      `json:"completed_at"`
}

// BatchJob tracks individual items in a batch
type BatchJob struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BatchID       uint      `json:"batch_id"`
	Batch         BatchOperation `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
	ItemID        uint      `json:"item_id"`
	ItemType      string    `json:"item_type"`
	Status        string    `json:"status"` // "pending", "processing", "completed", "failed"
	ErrorMessage  string    `json:"error_message"`
	Result        string    `json:"result"`
	StartedAt     time.Time `json:"started_at"`
	CompletedAt   time.Time `json:"completed_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateBatchOperationHandler creates a new batch operation
func (s *Server) createBatchOperationHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name       string      `json:"name"`
		Type       string      `json:"type"`
		TotalItems int         `json:"total_items"`
		Parameters interface{} `json:"parameters"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	paramJSON, _ := json.Marshal(req.Parameters)
	batch := BatchOperation{
		Name:       req.Name,
		Type:       req.Type,
		Status:     "pending",
		TotalItems: req.TotalItems,
		Parameters: paramJSON,
	}

	if err := s.db.Create(&batch).Error; err != nil {
		http.Error(w, "Failed to create batch operation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batch)
}

// GetBatchOperationsHandler retrieves batch operations
func (s *Server) getBatchOperationsHandler(w http.ResponseWriter, r *http.Request) {
	var batches []BatchOperation
	s.db.Order("created_at DESC").Limit(50).Find(&batches)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batches)
}

// GetBatchJobsHandler retrieves jobs in a batch
func (s *Server) getBatchJobsHandler(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchID")

	var jobs []BatchJob
	s.db.Where("batch_id = ?", batchID).Find(&jobs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// AddBatchJobsHandler adds items to a batch operation
func (s *Server) addBatchJobsHandler(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchID")

	var req struct {
		Items []struct {
			ItemID   uint   `json:"item_id"`
			ItemType string `json:"item_type"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var jobs []BatchJob
	for _, item := range req.Items {
		jobs = append(jobs, BatchJob{
			BatchID:  uint(0), // Parse from batchID
			ItemID:   item.ItemID,
			ItemType: item.ItemType,
			Status:   "pending",
		})
	}

	if err := s.db.CreateInBatches(jobs, 100).Error; err != nil {
		http.Error(w, "Failed to add batch jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"added": len(jobs)})
}

// StartBatchOperationHandler starts processing a batch
func (s *Server) startBatchOperationHandler(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchID")

	var batch BatchOperation
	if err := s.db.Where("id = ?", batchID).First(&batch).Error; err != nil {
		http.Error(w, "Batch not found", http.StatusNotFound)
		return
	}

	batch.Status = "in_progress"
	s.db.Save(&batch)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batch)
}

// GetBatchProgressHandler retrieves batch progress
func (s *Server) getBatchProgressHandler(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchID")

	var batch BatchOperation
	if err := s.db.Where("id = ?", batchID).First(&batch).Error; err != nil {
		http.Error(w, "Batch not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"batch_id":        batch.ID,
		"status":          batch.Status,
		"total":           batch.TotalItems,
		"processed":       batch.ProcessedItems,
		"failed":          batch.FailedItems,
		"progress_percent": float64(batch.ProcessedItems) / float64(batch.TotalItems) * 100,
	})
}

// CancelBatchOperationHandler cancels a batch operation
func (s *Server) cancelBatchOperationHandler(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchID")

	var batch BatchOperation
	if err := s.db.Where("id = ?", batchID).First(&batch).Error; err != nil {
		http.Error(w, "Batch not found", http.StatusNotFound)
		return
	}

	batch.Status = "cancelled"
	s.db.Save(&batch)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Batch cancelled"})
}

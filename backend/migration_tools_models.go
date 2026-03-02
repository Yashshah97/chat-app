package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Migration tracks data migration jobs
type Migration struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `json:"name"`
	Type          string         `json:"type"` // "import", "export", "sync"
	Source        string         `json:"source"`
	Destination   string         `json:"destination"`
	Status        string         `json:"status"` // "pending", "in_progress", "completed", "failed"
	TotalRecords  int            `json:"total_records"`
	MigratedCount int            `json:"migrated_count"`
	FailedCount   int            `json:"failed_count"`
	ErrorLog      datatypes.JSON `gorm:"type:jsonb" json:"error_log"`
	StartedAt     time.Time      `json:"started_at"`
	CompletedAt   time.Time      `json:"completed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// MigrationMap defines field mappings for migrations
type MigrationMap struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	MigrationID uint           `json:"migration_id"`
	Migration   Migration      `gorm:"foreignKey:MigrationID" json:"migration,omitempty"`
	SourceField string         `json:"source_field"`
	DestField   string         `json:"dest_field"`
	Transform   string         `json:"transform"` // transformation rule
	Config      datatypes.JSON `gorm:"type:jsonb" json:"config"`
	CreatedAt   time.Time      `json:"created_at"`
}

// MigrationTemplate provides pre-built migration configs
type MigrationTemplate struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	SourceSystem    string         `json:"source_system"`
	DestSystem      string         `json:"dest_system"`
	FieldMappings   datatypes.JSON `gorm:"type:jsonb" json:"field_mappings"`
	CreatedAt       time.Time      `json:"created_at"`
}

// CreateMigrationHandler creates a new migration job
func (s *Server) createMigrationHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Source      string `json:"source"`
		Destination string `json:"destination"`
		TotalRecords int    `json:"total_records"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	migration := Migration{
		Name:           req.Name,
		Type:           req.Type,
		Source:         req.Source,
		Destination:    req.Destination,
		TotalRecords:   req.TotalRecords,
		Status:         "pending",
	}

	if err := s.db.Create(&migration).Error; err != nil {
		http.Error(w, "Failed to create migration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(migration)
}

// GetMigrationsHandler retrieves all migrations
func (s *Server) getMigrationsHandler(w http.ResponseWriter, r *http.Request) {
	var migrations []Migration
	s.db.Order("created_at DESC").Limit(50).Find(&migrations)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(migrations)
}

// GetMigrationHandler retrieves migration details
func (s *Server) getMigrationHandler(w http.ResponseWriter, r *http.Request) {
	migrationID := chi.URLParam(r, "migrationID")

	var migration Migration
	if err := s.db.Where("id = ?", migrationID).First(&migration).Error; err != nil {
		http.Error(w, "Migration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(migration)
}

// StartMigrationHandler starts a migration job
func (s *Server) startMigrationHandler(w http.ResponseWriter, r *http.Request) {
	migrationID := chi.URLParam(r, "migrationID")

	var migration Migration
	if err := s.db.Where("id = ?", migrationID).First(&migration).Error; err != nil {
		http.Error(w, "Migration not found", http.StatusNotFound)
		return
	}

	migration.Status = "in_progress"
	migration.StartedAt = time.Now()
	s.db.Save(&migration)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(migration)
}

// GetMigrationProgressHandler retrieves migration progress
func (s *Server) getMigrationProgressHandler(w http.ResponseWriter, r *http.Request) {
	migrationID := chi.URLParam(r, "migrationID")

	var migration Migration
	if err := s.db.Where("id = ?", migrationID).First(&migration).Error; err != nil {
		http.Error(w, "Migration not found", http.StatusNotFound)
		return
	}

	progress := float64(migration.MigratedCount) / float64(migration.TotalRecords) * 100
	if migration.TotalRecords == 0 {
		progress = 0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"migration_id": migration.ID,
		"status":       migration.Status,
		"total":        migration.TotalRecords,
		"migrated":     migration.MigratedCount,
		"failed":       migration.FailedCount,
		"progress":     progress,
	})
}

// CreateMigrationMapHandler creates field mappings for migration
func (s *Server) createMigrationMapHandler(w http.ResponseWriter, r *http.Request) {
	migrationID := chi.URLParam(r, "migrationID")

	var req struct {
		SourceField string      `json:"source_field"`
		DestField   string      `json:"dest_field"`
		Transform   string      `json:"transform"`
		Config      interface{} `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	configJSON, _ := json.Marshal(req.Config)
	mapRecord := MigrationMap{
		MigrationID: uint(0), // Parse from migrationID
		SourceField: req.SourceField,
		DestField:   req.DestField,
		Transform:   req.Transform,
		Config:      configJSON,
	}

	if err := s.db.Create(&mapRecord).Error; err != nil {
		http.Error(w, "Failed to create migration map", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapRecord)
}

// GetMigrationTemplatesHandler retrieves available migration templates
func (s *Server) getMigrationTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	var templates []MigrationTemplate
	s.db.Find(&templates)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

// ValidateMigrationHandler validates a migration configuration
func (s *Server) validateMigrationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		TotalRecords int    `json:"total_records"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate configuration
	isValid := true
	var errors []string

	if req.Source == "" {
		isValid = false
		errors = append(errors, "Source system is required")
	}

	if req.Destination == "" {
		isValid = false
		errors = append(errors, "Destination system is required")
	}

	if req.TotalRecords <= 0 {
		isValid = false
		errors = append(errors, "Total records must be greater than 0")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":  isValid,
		"errors": errors,
	})
}

// GetMigrationStatsHandler retrieves migration statistics
func (s *Server) getMigrationStatsHandler(w http.ResponseWriter, r *http.Request) {
	var totalCount int64
	var completedCount int64
	var failedCount int64

	s.db.Model(&Migration{}).Count(&totalCount)
	s.db.Model(&Migration{}).Where("status = ?", "completed").Count(&completedCount)
	s.db.Model(&Migration{}).Where("status = ?", "failed").Count(&failedCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":     totalCount,
		"completed": completedCount,
		"failed":    failedCount,
		"success_rate": float64(completedCount) / float64(totalCount),
	})
}

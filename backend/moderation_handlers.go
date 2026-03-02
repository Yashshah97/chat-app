package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/admin/reports - Create a user report
func (s *Server) createReportHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		ReporterID uint   `json:"reporter_id"`
		TargetID   uint   `json:"target_id"`
		Reason     string `json:"reason"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	report := UserReport{
		ReporterID: reqBody.ReporterID,
		TargetID:   reqBody.TargetID,
		Reason:     reqBody.Reason,
		Status:     "pending",
	}
	
	result := s.db.Create(&report)
	if result.Error != nil {
		http.Error(w, "Failed to create report", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

// GET /api/admin/reports - List all reports
func (s *Server) listReportsHandler(w http.ResponseWriter, r *http.Request) {
	var reports []UserReport
	result := s.db.Preload("Reporter").Preload("Target").Find(&reports)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch reports", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
}

// GET /api/admin/reports/{id} - Get report details
func (s *Server) getReportDetailsHandler(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(reportID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}
	
	var report UserReport
	result := s.db.Preload("Reporter").Preload("Target").First(&report, id)
	
	if result.Error != nil {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// POST /api/admin/reports/{id}/resolve - Resolve a report
func (s *Server) resolveReportHandler(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(reportID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}
	
	var reqBody struct {
		Status     string `json:"status"` // resolved, dismissed, investigating
		Resolution string `json:"resolution"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result := s.db.Model(&UserReport{}, id).Updates(map[string]interface{}{
		"status":     reqBody.Status,
		"resolution": reqBody.Resolution,
	})
	
	if result.Error != nil {
		http.Error(w, "Failed to update report", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Report resolved successfully"})
}

// POST /api/admin/actions - Log admin action
func (s *Server) logAdminActionHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		AdminID    uint   `json:"admin_id"`
		Action     string `json:"action"`
		TargetID   uint   `json:"target_id"`
		TargetType string `json:"target_type"`
		Reason     string `json:"reason"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	action := AdminAction{
		AdminID:    reqBody.AdminID,
		Action:     reqBody.Action,
		TargetID:   reqBody.TargetID,
		TargetType: reqBody.TargetType,
		Reason:     reqBody.Reason,
	}
	
	result := s.db.Create(&action)
	if result.Error != nil {
		http.Error(w, "Failed to log action", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(action)
}

// GET /api/admin/actions - List all admin actions
func (s *Server) listAdminActionsHandler(w http.ResponseWriter, r *http.Request) {
	var actions []AdminAction
	result := s.db.Preload("Admin").Find(&actions)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch actions", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actions)
}

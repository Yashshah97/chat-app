package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/notifications/templates - Create a notification template
func (s *Server) createNotificationTemplateHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody NotificationTemplate
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result := s.db.Create(&reqBody)
	if result.Error != nil {
		http.Error(w, "Failed to create template", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reqBody)
}

// GET /api/notifications/templates - List all templates
func (s *Server) listNotificationTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"
	
	var templates []NotificationTemplate
	query := s.db
	
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	
	result := query.Order("name").Find(&templates)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch templates", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(templates)
}

// GET /api/notifications/templates/{id} - Get template details
func (s *Server) getNotificationTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(templateID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}
	
	var template NotificationTemplate
	result := s.db.First(&template, id)
	
	if result.Error != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(template)
}

// PUT /api/notifications/templates/{id} - Update template
func (s *Server) updateNotificationTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(templateID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}
	
	var reqBody NotificationTemplate
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result := s.db.Model(&NotificationTemplate{}, id).Updates(reqBody)
	
	if result.Error != nil {
		http.Error(w, "Failed to update template", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Template updated successfully"})
}

// DELETE /api/notifications/templates/{id} - Delete template
func (s *Server) deleteNotificationTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(templateID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Delete(&NotificationTemplate{}, id)
	
	if result.Error != nil {
		http.Error(w, "Failed to delete template", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Template deleted successfully"})
}

// POST /api/notifications/from-template - Create notification from template
func (s *Server) createNotificationFromTemplateHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		TemplateID uint            `json:"template_id"`
		UserID     uint            `json:"user_id"`
		Variables  map[string]string `json:"variables"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	var template NotificationTemplate
	if result := s.db.First(&template, reqBody.TemplateID); result.Error != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
	
	// In a real implementation, you'd interpolate variables into template
	notification := Notification{
		UserID:   reqBody.UserID,
		Type:     template.Type,
		Title:    template.Subject,
		Content:  template.Body,
		Priority: "normal",
		IsRead:   false,
	}
	
	result := s.db.Create(&notification)
	if result.Error != nil {
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

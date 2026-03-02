package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// POST /api/notifications/{id}/deliver - Send notification to delivery channels
func (s *Server) sendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	notifID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(notifID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	
	var reqBody struct {
		Channels []string `json:"channels"` // push, email, sms, in_app
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Get notification
	var notification Notification
	if result := s.db.First(&notification, id); result.Error != nil {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}
	
	// Create delivery records
	var deliveries []NotificationDelivery
	for _, channel := range reqBody.Channels {
		delivery := NotificationDelivery{
			NotificationID: uint(id),
			Channel:        channel,
			Status:         "pending",
			RetryCount:     0,
		}
		deliveries = append(deliveries, delivery)
	}
	
	result := s.db.CreateInBatches(deliveries, len(deliveries))
	if result.Error != nil {
		http.Error(w, "Failed to create deliveries", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Notification queued for delivery",
		"count": len(deliveries),
		"deliveries": deliveries,
	})
}

// GET /api/notifications/{id}/delivery-status - Get delivery status
func (s *Server) getDeliveryStatusHandler(w http.ResponseWriter, r *http.Request) {
	notifID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(notifID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	
	var deliveries []NotificationDelivery
	result := s.db.Where("notification_id = ?", id).Find(&deliveries)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch delivery status", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deliveries)
}

// PUT /api/notifications/delivery/{id}/mark-delivered - Mark as delivered
func (s *Server) markDeliveredHandler(w http.ResponseWriter, r *http.Request) {
	deliveryID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(deliveryID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}
	
	now := time.Now()
	result := s.db.Model(&NotificationDelivery{}, id).Updates(map[string]interface{}{
		"status": "delivered",
		"delivered_at": now,
	})
	
	if result.Error != nil {
		http.Error(w, "Failed to update delivery", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Delivery marked as complete"})
}

// POST /api/notifications/delivery/{id}/retry - Retry failed delivery
func (s *Server) retryDeliveryHandler(w http.ResponseWriter, r *http.Request) {
	deliveryID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(deliveryID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Model(&NotificationDelivery{}, id).Updates(map[string]interface{}{
		"status": "pending",
		"retry_count": s.db.Model(&NotificationDelivery{}).Session(&gormSession{}).Select("retry_count").Where("id = ?", id),
	})
	
	if result.Error != nil {
		http.Error(w, "Failed to retry delivery", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Delivery queued for retry"})
}

// GET /api/notifications/delivery/stats - Get delivery statistics
func (s *Server) getDeliveryStatsHandler(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		Total     int64
		Sent      int64
		Failed    int64
		Delivered int64
		Pending   int64
	}
	
	s.db.Model(&NotificationDelivery{}).Count(&stats.Total)
	s.db.Model(&NotificationDelivery{}).Where("status = ?", "sent").Count(&stats.Sent)
	s.db.Model(&NotificationDelivery{}).Where("status = ?", "failed").Count(&stats.Failed)
	s.db.Model(&NotificationDelivery{}).Where("status = ?", "delivered").Count(&stats.Delivered)
	s.db.Model(&NotificationDelivery{}).Where("status = ?", "pending").Count(&stats.Pending)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// Stub for gorm session interface
type gormSession struct{}

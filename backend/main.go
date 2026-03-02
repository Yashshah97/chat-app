package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
	router *chi.Mux
}

func main() {
	// Initialize upload directory
	if err := initializeUploadDirectory(); err != nil {
		log.Fatalf("Failed to initialize upload directory: %v", err)
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&User{}, 
		&Message{}, 
		&Chat{},
		&AdminUser{},
		&AdminAction{},
		&UserReport{},
		&ChatAnalytics{},
		&UserAnalytics{},
		&SystemAnalytics{},
		&UserPresence{},
		&PresenceHistory{},
		&ChatSettings{},
		&UserChatPreference{},
		&NotificationPreference{},
		&PinnedMessage{},
		&BlockedUser{},
		&MutedUser{},
		&ForwardedMessage{},
		&Notification{},
		&NotificationDelivery{},
		&NotificationTemplate{},
		&MessageEdit{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize server
	server := &Server{
		db: db,
		router: chi.NewRouter(),
	}

	// Setup routes
	server.setupRoutes()
	server.reactionHandlers()
	server.typingHandlers()
	server.readReceiptHandlers()

	// Start server
	port := ":8080"
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, server.router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(corsMiddleware)

	// Auth routes
	s.router.Post("/api/auth/register", s.registerHandler)
	s.router.Post("/api/auth/login", s.loginHandler)

	// Chat routes
	s.router.With(authMiddleware).Get("/api/chats", s.getChatsHandler)
	s.router.With(authMiddleware).Post("/api/chats", s.createChatHandler)
	s.router.With(authMiddleware).Post("/api/chats/group", s.createGroupChatHandler)
	s.router.With(authMiddleware).Get("/api/chats/{id}/messages", s.getChatMessagesHandler)
	
	// File upload routes
	s.router.With(authMiddleware).Post("/api/files/upload", uploadFileHandler)
	s.router.Get("/api/files/download/{filename}", downloadFileHandler)
	s.router.With(authMiddleware).Delete("/api/files/delete/{filename}", deleteFileHandler)
	s.router.Get("/api/files/list", listFilesHandler)

	// WebSocket
	s.router.With(authMiddleware).Get("/ws/chat/{id}", s.websocketHandler)
	
	// Admin routes - User management
	s.router.With(authMiddleware).Get("/api/admin/users", s.listUsersHandler)
	s.router.With(authMiddleware).Get("/api/admin/users/{id}", s.getUserDetailsHandler)
	s.router.With(authMiddleware).Post("/api/admin/users/{id}/suspend", s.suspendUserHandler)
	s.router.With(authMiddleware).Post("/api/admin/users/{id}/unsuspend", s.unsuspendUserHandler)
	s.router.With(authMiddleware).Delete("/api/admin/users/{id}", s.deleteUserHandler)
	
	// Admin routes - Chat management
	s.router.With(authMiddleware).Get("/api/admin/chats", s.listChatsHandler)
	s.router.With(authMiddleware).Get("/api/admin/chats/{id}", s.getChatDetailsHandler)
	s.router.With(authMiddleware).Delete("/api/admin/chats/{id}", s.deleteChatHandler)
	s.router.With(authMiddleware).Post("/api/admin/chats/{id}/remove-member/{memberID}", s.removeChatMemberHandler)
	s.router.With(authMiddleware).Post("/api/admin/chats/{id}/mute", s.muteChatHandler)
	
	// Admin routes - Moderation & Reporting
	s.router.Post("/api/reports", s.createReportHandler)
	s.router.With(authMiddleware).Get("/api/admin/reports", s.listReportsHandler)
	s.router.With(authMiddleware).Get("/api/admin/reports/{id}", s.getReportDetailsHandler)
	s.router.With(authMiddleware).Post("/api/admin/reports/{id}/resolve", s.resolveReportHandler)
	s.router.With(authMiddleware).Post("/api/admin/actions", s.logAdminActionHandler)
	s.router.With(authMiddleware).Get("/api/admin/actions", s.listAdminActionsHandler)
	
	// Analytics routes
	s.router.With(authMiddleware).Get("/api/analytics/system", s.getSystemAnalyticsHandler)
	s.router.With(authMiddleware).Get("/api/analytics/chat/{id}", s.getChatAnalyticsHandler)
	s.router.With(authMiddleware).Get("/api/analytics/user/{id}", s.getUserAnalyticsHandler)
	s.router.With(authMiddleware).Post("/api/analytics/compute", s.computeAnalyticsHandler)
	s.router.With(authMiddleware).Get("/api/analytics/dashboard", s.getAnalyticsDashboardHandler)
	
	// Presence/Online status routes
	s.router.With(authMiddleware).Post("/api/presence/update", s.updatePresenceHandler)
	s.router.With(authMiddleware).Get("/api/presence/user/{id}", s.getUserPresenceHandler)
	s.router.With(authMiddleware).Get("/api/presence/chat/{id}/members", s.getChatMembersPresenceHandler)
	s.router.Get("/api/presence/online", s.getOnlineUsersCountHandler)
	s.router.With(authMiddleware).Post("/api/presence/set-away", s.setUserAwayHandler)
	s.router.With(authMiddleware).Post("/api/presence/history", s.logPresenceHistoryHandler)
	
	// Search routes
	s.router.With(authMiddleware).Post("/api/search/messages", s.searchMessagesHandler)
	s.router.With(authMiddleware).Get("/api/search/messages/chat/{id}", s.searchChatMessagesHandler)
	s.router.With(authMiddleware).Get("/api/search/users", s.searchUsersHandler)
	s.router.With(authMiddleware).Get("/api/search/chats", s.searchChatsHandler)
	s.router.With(authMiddleware).Post("/api/search/advanced", s.advancedSearchHandler)
	s.router.Get("/api/search/trending", s.getTrendingHandler)
	s.router.With(authMiddleware).Get("/api/search/history/{userID}", s.getSearchHistoryHandler)
	
	// Chat settings routes
	s.router.With(authMiddleware).Get("/api/chats/{id}/settings", s.getChatSettingsHandler)
	s.router.With(authMiddleware).Put("/api/chats/{id}/settings", s.updateChatSettingsHandler)
	s.router.With(authMiddleware).Get("/api/users/{id}/preferences/{chatID}", s.getUserChatPreferencesHandler)
	s.router.With(authMiddleware).Put("/api/users/{id}/preferences/{chatID}", s.updateUserChatPreferencesHandler)
	s.router.With(authMiddleware).Get("/api/users/{id}/notifications", s.getNotificationPreferencesHandler)
	s.router.With(authMiddleware).Put("/api/users/{id}/notifications", s.updateNotificationPreferencesHandler)
	s.router.With(authMiddleware).Post("/api/chats/{id}/mute/{userID}", s.muteChatForUserHandler)
	
	// Message pinning routes
	s.router.With(authMiddleware).Post("/api/messages/{id}/pin", s.pinMessageHandler)
	s.router.With(authMiddleware).Delete("/api/messages/{id}/pin", s.unpinMessageHandler)
	s.router.With(authMiddleware).Get("/api/chats/{id}/pinned", s.getPinnedMessagesHandler)
	s.router.With(authMiddleware).Get("/api/messages/{id}/pin-status", s.checkPinStatusHandler)
	
	// User blocking and muting routes
	s.router.With(authMiddleware).Post("/api/users/{id}/block/{targetID}", s.blockUserHandler)
	s.router.With(authMiddleware).Delete("/api/users/{id}/block/{targetID}", s.unblockUserHandler)
	s.router.With(authMiddleware).Get("/api/users/{id}/blocked", s.getBlockedUsersHandler)
	s.router.With(authMiddleware).Post("/api/users/{id}/mute/{targetID}", s.muteUserHandler)
	s.router.With(authMiddleware).Delete("/api/users/{id}/mute/{targetID}", s.unmuteUserHandler)
	s.router.With(authMiddleware).Get("/api/users/{id}/muted", s.getMutedUsersHandler)
	s.router.With(authMiddleware).Get("/api/users/{id}/is-blocked-by/{targetID}", s.checkBlockStatusHandler)
	
	// Message forwarding routes
	s.router.With(authMiddleware).Post("/api/messages/{id}/forward", s.forwardMessageHandler)
	s.router.With(authMiddleware).Get("/api/chats/{id}/forwarded", s.getForwardedMessagesHandler)
	s.router.With(authMiddleware).Post("/api/messages/{id}/forward-to-multiple", s.forwardToMultipleHandler)
	
	// Notification routes
	s.router.With(authMiddleware).Post("/api/notifications", s.createNotificationHandler)
	s.router.With(authMiddleware).Get("/api/users/{id}/notifications", s.getUserNotificationsHandler)
	s.router.With(authMiddleware).Put("/api/notifications/{id}/read", s.markNotificationReadHandler)
	s.router.With(authMiddleware).Post("/api/users/{id}/notifications/mark-all-read", s.markAllNotificationsReadHandler)
	s.router.With(authMiddleware).Delete("/api/notifications/{id}", s.deleteNotificationHandler)
	s.router.With(authMiddleware).Get("/api/users/{id}/notifications/unread-count", s.getUnreadNotificationCountHandler)
	
	// Notification delivery routes
	s.router.With(authMiddleware).Post("/api/notifications/{id}/deliver", s.sendNotificationHandler)
	s.router.With(authMiddleware).Get("/api/notifications/{id}/delivery-status", s.getDeliveryStatusHandler)
	s.router.With(authMiddleware).Put("/api/notifications/delivery/{id}/mark-delivered", s.markDeliveredHandler)
	s.router.With(authMiddleware).Post("/api/notifications/delivery/{id}/retry", s.retryDeliveryHandler)
	s.router.Get("/api/notifications/delivery/stats", s.getDeliveryStatsHandler)
	
	// Notification template routes
	s.router.With(authMiddleware).Post("/api/notifications/templates", s.createNotificationTemplateHandler)
	s.router.With(authMiddleware).Get("/api/notifications/templates", s.listNotificationTemplatesHandler)
	s.router.With(authMiddleware).Get("/api/notifications/templates/{id}", s.getNotificationTemplateHandler)
	s.router.With(authMiddleware).Put("/api/notifications/templates/{id}", s.updateNotificationTemplateHandler)
	s.router.With(authMiddleware).Delete("/api/notifications/templates/{id}", s.deleteNotificationTemplateHandler)
	s.router.With(authMiddleware).Post("/api/notifications/from-template", s.createNotificationFromTemplateHandler)
	
	// Message editing routes
	s.router.With(authMiddleware).Put("/api/messages/{id}", s.editMessageHandler)
	s.router.With(authMiddleware).Get("/api/messages/{id}/history", s.getMessageHistoryHandler)
	s.router.With(authMiddleware).Get("/api/messages/{id}/edit-count", s.getEditCountHandler)
	s.router.With(authMiddleware).Delete("/api/message-edits/{id}", s.deleteEditHistoryHandler)
	s.router.With(authMiddleware).Get("/api/chats/{id}/edited-messages", s.getEditedMessagesHandler)

	// Health check
	s.router.Get("/health", s.healthHandler)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}
		
		// Verify token (simplified for now)
		if len(token) < 7 || token[:7] != "Bearer " {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Handler stubs
func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"jwt-token","username":"user"}`))
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"jwt-token","username":"user"}`))
}

func (s *Server) getChatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func (s *Server) createChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"id":1,"name":"Chat","type":"private"}`))
}

func (s *Server) createGroupChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"id":1,"name":"Group","type":"group"}`))
}

func (s *Server) getChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"websocket connection established"}`))
}

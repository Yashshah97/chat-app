package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/permissions - Create a new permission
func (s *Server) createPermissionHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		Category     string `json:"category"`
		ResourceType string `json:"resource_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	permission := Permission{
		Name:         reqBody.Name,
		Description:  reqBody.Description,
		Category:     reqBody.Category,
		ResourceType: reqBody.ResourceType,
	}

	result := s.db.Create(&permission)
	if result.Error != nil {
		http.Error(w, "Failed to create permission", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(permission)
}

// GET /api/permissions - List all permissions
func (s *Server) listPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	
	query := s.db
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var permissions []Permission
	result := query.Find(&permissions)
	if result.Error != nil {
		http.Error(w, "Failed to fetch permissions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permissions)
}

// POST /api/roles - Create a new role
func (s *Server) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Level       int    `json:"level"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role := Role{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Level:       reqBody.Level,
	}

	result := s.db.Create(&role)
	if result.Error != nil {
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

// GET /api/roles - List all roles
func (s *Server) listRolesHandler(w http.ResponseWriter, r *http.Request) {
	var roles []Role
	result := s.db.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}

// POST /api/roles/{id}/permissions/{permID} - Add permission to role
func (s *Server) addPermissionToRoleHandler(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	permID := chi.URLParam(r, "permID")

	rid, _ := strconv.ParseUint(roleID, 10, 32)
	pid, err := strconv.ParseUint(permID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid IDs", http.StatusBadRequest)
		return
	}

	var role Role
	if err := s.db.First(&role, rid).Error; err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	result := s.db.Model(&role).Association("Permissions").Append(&Permission{}, uint(pid))
	if result != nil {
		http.Error(w, "Failed to add permission", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// POST /api/users/{id}/roles/{roleID} - Assign role to user
func (s *Server) assignRoleToUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	roleID := chi.URLParam(r, "roleID")

	uid, _ := strconv.ParseUint(userID, 10, 32)
	rid, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid IDs", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		ExpiresAt *string `json:"expires_at"`
	}
	json.NewDecoder(r.Body).Decode(&reqBody)

	grantedByID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	userRole := UserRole{
		UserID:    uint(uid),
		RoleID:    uint(rid),
		GrantedBy: uint(grantedByID),
		ExpiresAt: reqBody.ExpiresAt,
	}

	result := s.db.Create(&userRole)
	if result.Error != nil {
		http.Error(w, "Failed to assign role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRole)
}

// GET /api/users/{id}/roles - Get user's roles
func (s *Server) getUserRolesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var roles []UserRole
	result := s.db.
		Where("user_id = ?", uid).
		Preload("Role").
		Preload("Role.Permissions").
		Find(&roles)

	if result.Error != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}

// POST /api/chats/{id}/roles/{roleID}/users/{userID} - Assign role in chat
func (s *Server) assignChatRoleHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	roleID := chi.URLParam(r, "roleID")
	userID := chi.URLParam(r, "userID")

	cid, _ := strconv.ParseUint(chatID, 10, 32)
	rid, _ := strconv.ParseUint(roleID, 10, 32)
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid IDs", http.StatusBadRequest)
		return
	}

	grantedByID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	userRole := UserRole{
		UserID:    uint(uid),
		RoleID:    uint(rid),
		ChatID:    (*uint)(&cid),
		GrantedBy: uint(grantedByID),
	}

	result := s.db.Create(&userRole)
	if result.Error != nil {
		http.Error(w, "Failed to assign role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRole)
}

// DELETE /api/users/{id}/roles/{roleID} - Revoke role from user
func (s *Server) revokeRoleHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	roleID := chi.URLParam(r, "roleID")

	uid, _ := strconv.ParseUint(userID, 10, 32)
	rid, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid IDs", http.StatusBadRequest)
		return
	}

	result := s.db.
		Where("user_id = ? AND role_id = ?", uid, rid).
		Delete(&UserRole{})

	if result.Error != nil {
		http.Error(w, "Failed to revoke role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "revoked"})
}

// GET /api/permissions/audit - Get permission audit logs
func (s *Server) getPermissionAuditHandler(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var audits []PermissionAudit
	result := s.db.
		Order("created_at DESC").
		Limit(limit).
		Find(&audits)

	if result.Error != nil {
		http.Error(w, "Failed to fetch audits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(audits)
}

// GET /api/users/{id}/permissions - Get user's effective permissions
func (s *Server) getUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var permissions []Permission
	result := s.db.
		Table("permissions").
		Joins("INNER JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("INNER JOIN roles ON roles.id = role_permissions.role_id").
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", uid).
		Distinct("permissions.*").
		Find(&permissions)

	if result.Error != nil {
		http.Error(w, "Failed to fetch permissions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permissions)
}

package main

import "gorm.io/gorm"

// Permission represents a specific action permission
type Permission struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"` // admin, chat, message, file, etc
	ResourceType string `json:"resource_type"` // chat, message, user, etc
}

// Role represents a group of permissions
type Role struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Description string `json:"description"`
	IsSystem    bool   `gorm:"default:false" json:"is_system"` // Built-in roles: admin, moderator, member
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Level       int    `gorm:"default:0" json:"level"` // Higher number = more privileged
}

// UserRole represents roles assigned to users
type UserRole struct {
	gorm.Model
	UserID   uint   `gorm:"not null;index" json:"user_id"`
	User     *User  `gorm:"foreignKey:UserID" json:"-"`
	RoleID   uint   `gorm:"not null;index" json:"role_id"`
	Role     *Role  `gorm:"foreignKey:RoleID" json:"-"`
	ChatID   *uint  `json:"chat_id"` // NULL for system roles, specified for chat-specific roles
	Chat     *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	GrantedBy uint  `json:"granted_by"`
	ExpiresAt *string `json:"expires_at"` // Temporary roles
}

// ChatPermission represents permissions for specific actions in a chat
type ChatPermission struct {
	gorm.Model
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	RoleID      uint   `gorm:"not null;index" json:"role_id"`
	Role        *Role  `gorm:"foreignKey:RoleID" json:"-"`
	PermissionID uint  `gorm:"not null;index" json:"permission_id"`
	Permission  *Permission `gorm:"foreignKey:PermissionID" json:"-"`
	IsGranted   bool   `gorm:"default:true" json:"is_granted"` // Explicit deny/grant
	GrantedAt   string `json:"granted_at"`
	GrantedByID uint   `json:"granted_by_id"`
}

// PermissionAudit tracks permission changes
type PermissionAudit struct {
	gorm.Model
	UserID       uint   `gorm:"not null;index" json:"user_id"`
	User         *User  `gorm:"foreignKey:UserID" json:"-"`
	Action       string `json:"action"` // grant, revoke, assign_role, remove_role
	TargetUserID *uint  `json:"target_user_id"`
	TargetUser   *User  `gorm:"foreignKey:TargetUserID" json:"-"`
	ResourceType string `json:"resource_type"` // chat, permission, role
	ResourceID   uint   `json:"resource_id"`
	OldValue     string `json:"old_value"`
	NewValue     string `json:"new_value"`
	Reason       string `json:"reason"`
	IPAddress    string `json:"ip_address"`
}

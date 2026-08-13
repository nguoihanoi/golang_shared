package permissions

import (
	"time"
)

// PermissionType represents a permission group in the system
type PermissionType struct {
	ID        string              `bson:"_id,omitempty" json:"_id,omitempty"`
	Delete    int                 `bson:"delete" json:"delete"`
	CreatedAt time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at" json:"updated_at"`
	Name      map[string]string   `bson:"name" json:"name"`
	Values    map[string][]string `bson:"values" json:"values"`
	Order     int                 `bson:"order" json:"order"`
	AuthorId  string              `bson:"author_id" json:"author_id"`
}

// MiniPermissionType represents a permission group in the system
type MiniPermissionType struct {
	ID          string            `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        map[string]string `bson:"name" json:"name"`
	Permissions []MiniPermission  `bson:"permissions" json:"permissions"`
}

// MiniPermissionType represents a permission group in the system
type MiniPermission struct {
	ID   string            `bson:"_id,omitempty" json:"_id,omitempty"`
	Name map[string]string `bson:"name" json:"name"`
	Code string            `bson:"code" json:"code"`
}

// Permission represents a permission in the system
type Permission struct {
	ID               string            `bson:"_id,omitempty" json:"_id,omitempty"`
	Delete           int               `bson:"delete" json:"delete"`
	CreatedAt        time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time         `bson:"updated_at" json:"updated_at"`
	PermissionTypeID string            `bson:"type_id" json:"type_id"`
	Name             map[string]string `bson:"name" json:"name"`
	Code             string            `bson:"code" json:"code"`
	Order            int               `bson:"order" json:"order"`
	AuthorId         string            `bson:"author_id" json:"author_id"`
}

// Account type represents a account type in the system
type AccountType struct {
	ID          string            `bson:"_id,omitempty" json:"_id,omitempty"`
	Delete      int               `bson:"delete" json:"delete"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
	Name        map[string]string `bson:"name" json:"name"`
	Content     map[string]string `bson:"content" json:"content"`
	Permissions []string          `bson:"permissions" json:"permissions"`
	Status      int               `bson:"status" json:"status"`
	Order       int               `bson:"order" json:"order"`
	AuthorId    string            `bson:"author_id" json:"author_id"`
}
type MiniAccountType struct {
	ID   string            `bson:"_id,omitempty" json:"_id,omitempty"`
	Name map[string]string `bson:"name" json:"name"`
}

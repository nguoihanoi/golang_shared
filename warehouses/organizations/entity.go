package organizations

import (
	"time"
)

// Organization represents a organization in the system
type Organization struct {
	ID          string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Delete      int       `bson:"delete" json:"delete"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Avatar      string    `bson:"avatar" json:"avatar"`
	TotalMember int       `bson:"total_member" json:"total_member"`
	Status      int       `bson:"status" json:"status"`
	AuthorId    string    `bson:"author_id" json:"author_id"`
}

// Member represents a permission in the system
type Member struct {
	ID             string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Delete         int       `bson:"delete" json:"delete"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
	OrganizationID string    `bson:"organization_id" json:"organization_id"`
	CustomerID     string    `bson:"customer_id" json:"customer_id"`
	Host           int       `bson:"host" json:"host"`
	AuthorId       string    `bson:"author_id" json:"author_id"`
}

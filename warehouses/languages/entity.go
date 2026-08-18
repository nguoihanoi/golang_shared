package languages

import (
	"time"
)

// Customer represents a customer in the system
type Language struct {
	ID        string    `bson:"_id" json:"_id"`
	Delete    int       `bson:"delete" json:"delete"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Name      string    `bson:"name" json:"name"`
	Code      string    `bson:"code" json:"code"`
	Image     string    `bson:"image" json:"image"`
	Order     int       `bson:"order" json:"order"`
	Status    int       `bson:"status" json:"status"`
	AuthorId  string    `bson:"author_id" json:"author_id"`
}

type LanguageCode struct {
	ID        string            `bson:"_id" json:"_id"`
	Delete    int               `bson:"delete" json:"delete"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
	Name      string            `bson:"name" json:"name"`
	Value     map[string]string `bson:"value" json:"value"`
	Type      int               `bson:"type" json:"type"`
	GroupId   string            `bson:"group_id" json:"group_id"`
	AuthorId  string            `bson:"author_id" json:"author_id"`
}
type GroupCode struct {
	ID        string    `bson:"_id" json:"_id"`
	Delete    int       `bson:"delete" json:"delete"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Name      string    `bson:"name" json:"name"`
	Code      string    `bson:"code" json:"code"`
	AuthorId  string    `bson:"author_id" json:"author_id"`
}

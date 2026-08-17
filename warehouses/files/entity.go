package files

import (
	"time"
)

type File struct {
	ID            string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Delete        int       `bson:"delete" json:"delete"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
	Name          string    `bson:"name" json:"name"`
	Size          int64     `bson:"size" json:"size"`
	Width         int       `bson:"width" json:"width"`
	Height        int       `bson:"height" json:"height"`
	Duration      float32   `bson:"duration" json:"duration"`
	Type          string    `bson:"type" json:"type"`
	S3Key         string    `bson:"s3_key,omitempty" json:"s3_key,omitempty"`
	S3Bucket      string    `bson:"s3_bucket,omitempty" json:"s3_bucket,omitempty"`
	StoredLocally bool      `bson:"stored_locally" json:"stored_locally"`
	LocalPath     string    `bson:"local_path,omitempty" json:"local_path,omitempty"`
	AuthorId      string    `bson:"author_id" json:"author_id"`
	Description   string    `bson:"description,omitempty" json:"description,omitempty"`
	Tags          []string  `bson:"tags,omitempty" json:"tags,omitempty"`
}

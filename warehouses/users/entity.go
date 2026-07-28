package users

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID               string    `bson:"_id" json:"_id"`
	Delete           int       `bson:"delete" json:"delete"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
	AccountType      string    `bson:"account_type" json:"account_type"`
	Email            string    `bson:"email" json:"email"`
	LanguageCode     string    `bson:"language_code" json:"language_code"`
	FirstName        string    `bson:"first_name" json:"first_name"`
	LastName         string    `bson:"last_name" json:"last_name"`
	PasswordHash     string    `bson:"password_hash" json:"password_hash"` // Never send this to client
	Password         string    `bson:"password" json:"password"`           // Never send this to client
	ResetToken       string    `bson:"reset_token" json:"reset_token"`
	ResetTokenExp    time.Time `bson:"reset_token_exp" json:"reset_token_exp"`
	IsActive         bool      `bson:"is_active" json:"is_active"`
	ProfilePicture   string    `bson:"profile_picture" json:"profile_picture"`
	AuthProvider     string    `bson:"auth_provider" json:"auth_provider"`           // "google", "facebook", or empty for local auth
	ProviderID       string    `bson:"provider_id" json:"provider_id"`               // ID from the social provider
	TwoFactorEnabled bool      `bson:"two_factor_enabled" json:"two_factor_enabled"` // Whether 2FA is enabled
}

// UserMeta represents a user meta in the system
type UserMeta struct {
	ID              string    `bson:"_id" json:"_id"`
	TwoFactorSecret string    `bson:"two_factor_secret" json:"-"` // Secret key for 2FA (TOTP)
	BackupCodes     []string  `bson:"backup_codes" json:"-"`      // Backup codes for 2FA recovery
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
}

type UserToken struct {
	UserID       string `bson:"user_id" json:"user_id"`
	LanguageCode string `bson:"lang_code" json:"lang_code"`
	Time         int64  `bson:"time" json:"time"`
}

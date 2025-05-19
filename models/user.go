package models

import "time"

type UserType string
type UserPlan string

const (
	UserTypeUser  UserType = "user"
	UserTypeAdmin UserType = "admin"
)

const (
	UserPlanFree    UserPlan = "free"
	UserPlanPremium UserPlan = "premium"
)

type User struct {
	// User ID (MongoDB ObjectID as hex string)
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	// Full name of the user
	Name string `json:"name" bson:"name"`

	// Email address of the user
	Email string `json:"email" bson:"email"`

	// Password hash (not returned in responses)
	Password string `json:"password,omitempty" bson:"password"`

	// Timestamp when the user was created
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`

	// Timestamp when the user was last updated
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`

	// Subscription plan of the user (e.g. "free", "premium")
	Plan string `json:"plan,omitempty" bson:"plan,omitempty"`

	// User type (e.g. "user", "admin")
	Type string `json:"type,omitempty" bson:"type,omitempty"`

	// Google OAuth ID (if user signed up with Google)
	GoogleID string `json:"googleId,omitempty" bson:"googleId,omitempty"`
}

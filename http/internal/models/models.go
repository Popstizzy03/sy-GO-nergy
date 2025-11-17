package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Add methods for template formatting
func (u User) initial() string {
	if len(u.Name) > 0 {
		return string(u.name[0])
	}
	return "?"
}

func (u User) FormatCreatedAt() string {
	return u.CreatedAt.Format("Jan 2, 2006")
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

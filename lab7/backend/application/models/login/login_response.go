package login

import "time"

type LoginResponse struct {
	ID                int64     `json:"id"`
	FullName          string    `json:"full_name"`
	SocialNetworkLink string    `json:"social_network_link"`
	PhoneNumber       string    `json:"phone_number"`
	Email             string    `json:"email"`
	BirthDate         time.Time `json:"birth_date"`
	Role              string    `json:"role"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	GroupName         string    `json:"group_name,omitempty"`
	AccessToken       string    `json:"access_token"`
	RefreshToken      string    `json:"refresh_token"`
}

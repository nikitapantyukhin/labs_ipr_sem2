package create_user

import (
	"sport_platform/application/models/shared"
	"time"
)

type CreateUserResponse struct {
	ID                int64           `json:"id"`
	FullName          string          `json:"full_name"`
	SocialNetworkLink string          `json:"social_network_link"`
	PhoneNumber       string          `json:"phone_number"`
	Email             string          `json:"email"`
	BirthDate         time.Time       `json:"birth_date"`
	Role              shared.UserRole `json:"role"`
	GroupID           *int64          `json:"group_id"`
	GroupName         string          `json:"group_name"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	AccessToken       string          `json:"access_token"`
	RefreshToken      string          `json:"refresh_token"`
}

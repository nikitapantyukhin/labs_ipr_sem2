package create_user

import "time"

type CreateUserRequest struct {
	FullName          string    `json:"full_name" validate:"required,min=2,max=50"`
	SocialNetworkLink string    `json:"social_network_link" validate:"required,min=2,startswith=@"`
	PhoneNumber       string    `json:"phone_number" validate:"required,e164"`
	Email             string    `json:"email" validate:"required,email"`
	BirthDate         time.Time `json:"birth_date"`
	Password          string    `json:"password" mapper:"exclude" validate:"required,min=8"`
	GroupID           *int64    `json:"group_id"`
}

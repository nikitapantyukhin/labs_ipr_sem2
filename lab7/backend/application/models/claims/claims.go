package claims

import "sport_platform/application/models/shared"

type UserClaims struct {
	ID       int64           `json:"id"`
	FullName string          `json:"full_name"`
	Email    string          `json:"email"`
	Role     shared.UserRole `json:"role"`
}

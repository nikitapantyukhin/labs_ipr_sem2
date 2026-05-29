package delete_user

type DeleteUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

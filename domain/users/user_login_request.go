package users

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ClientID int64  `json:"client_id"`
}

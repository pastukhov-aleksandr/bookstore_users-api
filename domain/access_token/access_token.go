package access_token

type AccessTokenRequest struct {
	UserID   int64  `json:"user_id"`
	ClientID int64  `json:"client_id"`
	UuID     string `json:"uuid"`
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Permission   string `json:"permission"`
}

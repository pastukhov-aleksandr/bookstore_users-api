package captcha

type CaptchaRequest struct {
	Email    string `json:"email"`
	Pin string `json:"pin"`
}
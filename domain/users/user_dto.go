package users

import (
	"strings"

	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
)

const (
	StatusActive = "active"
)

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	DateCreated  string `json:"date_created"`
	Status       string `json:"status"`
	Password     string `json:"password"`
	Captcha      string `json:"captcha"`
	UuID         string `json:"uuid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Permission   string `json:"permission"`
	Color        string `json:"color"`
}

type Users []User

func GetNewUsers(userID int64) User {
	return User{
		ID: userID,
	}
}

func (user *User) Validate() rest_errors.RestErr {
	user.Name = strings.TrimSpace(user.Name)
	user.Phone = strings.TrimSpace(user.Phone)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Captcha = strings.TrimSpace(user.Captcha)
	if user.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}
	if user.Captcha == "" {
		return rest_errors.NewBadRequestError("invalid Captcha")
	}
	return nil
}

package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pastukhov-aleksandr/bookstore_users-api/domain/access_token"
	"github.com/pastukhov-aleksandr/bookstore_users-api/domain/captcha"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
)

type RestUsersRepository interface {
	Captcha(string, string) rest_errors.RestErr
	CreateAccessToken(int64, string, int64) (*access_token.AccessToken, rest_errors.RestErr)
	DeleteRefreshToken(int64, int64) rest_errors.RestErr
}

type AuthSuccess struct {
	/* variables */
}

type usersRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (s *usersRepository) Captcha(email string, pin string) rest_errors.RestErr {
	request := captcha.CaptchaRequest{
		Email: email,
		Pin:   pin,
	}
	// Create a Resty Client
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&AuthSuccess{}).
		Post("http://localhost:8082/captcha/validate")

	if err != nil || response.Body() == nil {
		return rest_errors.NewInternalServerError("invalid restclient response when trying to captcha", errors.New("restclient error"))
	}

	if response.StatusCode() != http.StatusOK {
		apiErr, err := rest_errors.NewRestErrorFromBytes(response.Body())
		if err != nil {
			return rest_errors.NewInternalServerError("invalid captcha code", err)
		}
		return apiErr
	}

	return nil
}

func (s *usersRepository) CreateAccessToken(UserID int64, UuID string, ClientID int64) (*access_token.AccessToken, rest_errors.RestErr) {
	request := access_token.AccessTokenRequest{
		UserID:   UserID,
		ClientID: ClientID,
		UuID:     UuID,
	}
	// Create a Resty Client
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&AuthSuccess{}).
		Post("http://localhost:8080/oauth/access_token")

	if err != nil || response.Body() == nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to access token", errors.New("restclient error"))
	}

	if response.StatusCode() != http.StatusCreated {
		apiErr, err := rest_errors.NewRestErrorFromBytes(response.Body())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid access token", err)
		}
		return nil, apiErr
	}

	var at access_token.AccessToken
	if err := json.Unmarshal(response.Body(), &at); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal access token response", errors.New("json parsing error"))
	}

	return &at, nil
}

func (s *usersRepository) DeleteRefreshToken(userID int64, clientID int64) rest_errors.RestErr {
	request := access_token.AccessTokenRequest{
		UserID:   userID,
		ClientID: clientID,
		UuID:     "",
	}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("http://localhost:8080/oauth/logout")

	if err != nil {
		return rest_errors.NewInternalServerError("invalid restclient response when trying to access token", errors.New("restclient error"))
	}
	return nil
}

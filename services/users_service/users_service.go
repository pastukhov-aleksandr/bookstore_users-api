package users_service

import (
	"github.com/AvraamMavridis/randomcolor"
	"github.com/pastukhov-aleksandr/bookstore_users-api/domain/access_token"
	"github.com/pastukhov-aleksandr/bookstore_users-api/domain/users"
	"github.com/pastukhov-aleksandr/bookstore_users-api/repositore/db"
	"github.com/pastukhov-aleksandr/bookstore_users-api/repositore/rest"
	"github.com/pastukhov-aleksandr/bookstore_users-api/utils/crypto_utils"
	"github.com/pastukhov-aleksandr/bookstore_users-api/utils/date_utils"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	GetUser(int64) (*users.User, rest_errors.RestErr)
	CreateUser(users.User) (*users.User, rest_errors.RestErr)
	//UpdateUser(bool, users.User) (*users.User, rest_errors.RestErr)
	//DeleteUser(int64) rest_errors.RestErr
	//SearchUser(string) (users.Users, rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*access_token.AccessToken, rest_errors.RestErr)
	Logout(int64, int64) rest_errors.RestErr
	Refresh(int64, int64) (*access_token.AccessToken, rest_errors.RestErr)
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetUser(userID int64) (*users.User, rest_errors.RestErr) {
	result := users.GetNewUsers(userID)
	if err := s.dbRepo.Get(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *service) CreateUser(user users.User) (*users.User, rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// captcha {
	if err := s.restUsersRepo.Captcha(user.Email, user.Captcha); err != nil {
		return nil, err
	}
	//}

	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDbFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	user.Color = randomcolor.GetRandomColorInHex()
	if err := s.dbRepo.Save(user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) LoginUser(request users.LoginRequest) (*access_token.AccessToken, rest_errors.RestErr) {
	dao := &users.User{
		Email:      request.Email,
		Password:   crypto_utils.GetMd5(request.Password),
		UuID:       uuid.NewV4().String(),
		Permission: "aouth",
	}

	if err := s.dbRepo.FindByEmailAndPassword(dao); err != nil {
		return nil, err
	}

	// creating acsess token {
	at, err := s.restUsersRepo.CreateAccessToken(dao.ID, dao.UuID, request.ClientID)
	if err != nil {
		return nil, err
	}
	at.Permission = dao.Permission
	//}
	return at, nil
}

func (s *service) Logout(userID int64, clientID int64) rest_errors.RestErr {
	err := s.restUsersRepo.DeleteRefreshToken(userID, clientID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Refresh(userID int64, clientID int64) (*access_token.AccessToken, rest_errors.RestErr) {
	// delete old refresh tokens
	err := s.restUsersRepo.DeleteRefreshToken(userID, clientID)
	if err != nil {
		return nil, err
	}

	// creating new acsess token {
	at, err := s.restUsersRepo.CreateAccessToken(userID, uuid.NewV4().String(), clientID)
	if err != nil {
		return nil, err
	}
	at.Permission = "aouth"
	//}
	return at, nil
}

// func (s *service) UpdateUser(isPartial bool, user users.User) (*users.User, rest_errors.RestErr) {
// 	current, err := s.GetUser(user.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if isPartial {
// 		if user.Name != "" {
// 			current.Name = user.Name
// 		}

// 		if user.Email != "" {
// 			current.Email = user.Email
// 		}
// 	} else {
// 		current.Name = user.Name
// 		current.Email = user.Email
// 	}

// 	if err := s.dbRepo.Update(current); err != nil {
// 		return nil, err
// 	}

// 	return current, nil
// }

// func (s *service) DeleteUser(userID int64) rest_errors.RestErr {
// 	user := users.GetNewUsers(userID)
// 	return s.dbRepo.Delete(user)
// }

// func (s *service) SearchUser(status string) (users.Users, rest_errors.RestErr) {
// 	dao := &users.User{}
// 	return s.dbRepo.FindByStatus(status)
// }

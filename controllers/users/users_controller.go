package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pastukhov-aleksandr/bookstore_users-api/domain/users"
	"github.com/pastukhov-aleksandr/bookstore_users-api/services/users_service"
	"github.com/pastukhov-aleksandr/bookstore_users-api/utils/oauth"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/secret_code"
)

type UsersHandler interface {
	Create(*gin.Context)
	GetInfo(*gin.Context)
	//Update(*gin.Context)
	//Delete(*gin.Context)
	//Search(*gin.Context)
	Login(*gin.Context)
	Logout(*gin.Context)
	Refresh(*gin.Context)
}

type usersHandler struct {
	service users_service.Service
}

func NewUsersHandler(service users_service.Service) UsersHandler {
	return &usersHandler{
		service: service,
	}
}

func getUserId(userIdParam string) (int64, rest_errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("invalid user id")
	}
	return userId, nil
}

func (handler *usersHandler) Create(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, saveErr := handler.service.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status(), saveErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (handler *usersHandler) GetInfo(c *gin.Context) {
	ad, err := oauth.AuthenticateRequest(c.Request, secret_code.Get_ACCESS_SECRET())
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// userId, idErr := getUserId(c.Param("user_id"))
	// if idErr != nil {
	// 	c.JSON(idErr.Status(), idErr)
	// 	return
	// }
	user, getErr := handler.service.GetUser(ad.UserID)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	// if oauth.GetCallerId(c.Request) == user.ID {
	// 	c.JSON(http.StatusOK, user.Marshall(false))
	// 	return
	// }

	c.JSON(http.StatusOK, user.Marshall(false))
}

func (handler *usersHandler) Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	at, err := handler.service.LoginUser(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, at)
}

// func (handler *usersHandler) Update(c *gin.Context) {
// 	userId, idErr := getUserId(c.Param("user_id"))
// 	if idErr != nil {
// 		c.JSON(idErr.Status(), idErr)
// 		return
// 	}

// 	var user users.User

// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		restErr := rest_errors.NewBadRequestError("invalid json body")
// 		c.JSON(restErr.Status(), restErr)
// 		return
// 	}

// 	user.ID = userId

// 	isPartial := c.Request.Method == http.MethodPatch

// 	result, err := handler.service.UpdateUser(isPartial, user)
// 	if err != nil {
// 		c.JSON(err.Status(), err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
// }

// func (handler *usersHandler) Delete(c *gin.Context) {
// 	userId, idErr := getUserId(c.Param("user_id"))
// 	if idErr != nil {
// 		c.JSON(idErr.Status(), idErr)
// 		return
// 	}

// 	if err := handler.service.DeleteUser(userId); err != nil {
// 		c.JSON(err.Status(), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
// }

// func (handler *usersHandler) Search(c *gin.Context) {
// 	status := c.Query("status")

// 	users, err := handler.service.SearchUser(status)
// 	if err != nil {
// 		c.JSON(err.Status(), err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
// }

func (handler *usersHandler) Logout(c *gin.Context) {
	ad, err := oauth.AuthenticateRequest(c.Request, secret_code.Get_REFRESH_SECRET())
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	err = handler.service.Logout(ad.UserID, ad.ClientID)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func (handler *usersHandler) Refresh(c *gin.Context) {
	// ad, err := oauth.AuthenticateRequest(c.Request, secret_code.Get_REFRESH_SECRET())
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, err)
	// 	return
	// }

	// var request users.LoginRequest

	// at, err := handler.service.LoginUser(request)
	// if err != nil {
	// 	c.JSON(err.Status(), err)
	// 	return
	// }
	// c.JSON(http.StatusOK, at)
}

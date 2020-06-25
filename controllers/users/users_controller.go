package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pastukhov-aleksandr/bookstore_aouth-go/oauth"
	"github.com/pastukhov-aleksandr/bookstore_users-api/domain/users"
	"github.com/pastukhov-aleksandr/bookstore_users-api/services/users_service"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
)

type UsersHandler interface {
	Create(*gin.Context)
	Get(*gin.Context)
	//Update(*gin.Context)
	//Delete(*gin.Context)
	//Search(*gin.Context)
	Login(*gin.Context)
	//Logouts(*gin.Context)
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

func (handler *usersHandler) Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}
	user, getErr := handler.service.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	if oauth.GetCallerId(c.Request) == user.ID {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}

	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func (handler *usersHandler) Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	user, err := handler.service.LoginUser(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(false))
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

// func (handler *usersHandler) Logout(c *gin.Context) {
// 	// TODO: logout
// }

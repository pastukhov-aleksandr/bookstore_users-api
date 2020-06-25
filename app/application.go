package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pastukhov-aleksandr/bookstore_users-api/controllers/ping"
	"github.com/pastukhov-aleksandr/bookstore_users-api/controllers/users"
	"github.com/pastukhov-aleksandr/bookstore_users-api/repositore/db"
	"github.com/pastukhov-aleksandr/bookstore_users-api/repositore/rest"
	"github.com/pastukhov-aleksandr/bookstore_users-api/services/users_service"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/logger"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := users.NewUsersHandler(
		users_service.NewService(rest.NewRestUsersRepository(), db.NewRepository()))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	router.Use(cors.New(config))

	router.GET("/ping", ping.Ping)

	router.POST("/api/users", atHandler.Create)
	// router.GET("/api/users/:user_id", atHandler.Get)
	// router.PUT("/api/users/:user_id", atHandler.Update)
	// router.PATCH("/api/users/:user_id", atHandler.Update)
	// router.DELETE("/api/users/:user_id", atHandler.Delete)
	// router.GET("internal/users/search", atHandler.Search)
	router.POST("/api/users/login", atHandler.Login)
	// router.POST("/api/users/logout", atHandler.Logout)

	logger.Info("about to start the application...")
	router.Run(":8081")
}

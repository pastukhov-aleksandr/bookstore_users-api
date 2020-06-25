package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/logger"
)

const (
	mysqlUsersUsername = "mysql_users_username"
	mysqlUsersPassword = "mysql_users_password"
	mysqlUsersHost     = "mysql_users_host"
	mysqlUsersSchema   = "mysql_users_schema"
)

var (
	Client *sql.DB
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var (
		username = getEnv(mysqlUsersUsername, "") //os.Getenv(mysqlUsersUsername)
		password = getEnv(mysqlUsersPassword, "")
		host     = getEnv(mysqlUsersHost, "")
		schema   = getEnv(mysqlUsersSchema, "")
	)

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4",
		username, password, host, schema,
	)

	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	mysql.SetLogger(logger.GetLogger())
	log.Println("database successfully configured")
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

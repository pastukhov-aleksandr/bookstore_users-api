package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	users_db "github.com/pastukhov-aleksandr/bookstore_users-api/client/mysql"
	"github.com/pastukhov-aleksandr/bookstore_users-api/domain/users"
	"github.com/pastukhov-aleksandr/bookstore_users-api/utils/mysql_utils"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/logger"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
)

const (
	queryInsertUser             = "INSERT INTO users(name, email, date_created, status, password, phone, color) VALUES(?, ?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser             = "UPDATE users SET name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindByStatus           = "SELECT id, name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, name, email, date_created, status, color FROM users WHERE email=? AND password=? AND status=?"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	Get(users.User) rest_errors.RestErr
	Save(users.User) rest_errors.RestErr
	Update(users.User) rest_errors.RestErr
	Delete(users.User) rest_errors.RestErr
	FindByStatus(string) ([]users.User, rest_errors.RestErr)
	FindByEmailAndPassword(*users.User) rest_errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) Get(user users.User) rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.ID)
	if getErr := result.Scan(&user.ID, &user.Name, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error when trying to  get user by id", getErr)
		return rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}

	return nil
}

func (r *dbRepository) Save(user users.User) rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.Name, user.Email, user.DateCreated, user.Status, user.Password, user.Phone, user.Color)
	if saveErr != nil {
		fmt.Println(saveErr)
		if driverErr, ok := saveErr.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 {
				return rest_errors.NewInternalServerError("User with this email already exists", errors.New("database error"))
			}
		}

		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error when tying to save user", errors.New("database error"))
	}

	user.ID = userId
	return nil
}

func (r *dbRepository) Update(user users.User) rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when tying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Email, user.ID)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	return nil
}

func (r *dbRepository) Delete(user users.User) rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when tying to delete user", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.ID); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when tying to delete user", errors.New("database error"))
	}
	return nil
}

func (r *dbRepository) FindByStatus(status string) ([]users.User, rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when tying to find by status user", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer rows.Close()

	result := make([]users.User, 0)
	for rows.Next() {
		var user users.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scan user row struct", err)
			return nil, rest_errors.NewInternalServerError("error when tying to find by status user", errors.New("database error"))
		}
		result = append(result, user)
	}
	if len(result) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no user matching status %s", status))
	}
	return result, nil
}

func (r *dbRepository) FindByEmailAndPassword(user *users.User) rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when tying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, users.StatusActive)
	if getErr := result.Scan(&user.ID, &user.Name, &user.Email, &user.DateCreated, &user.Status, &user.Color); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("error when tying to find user", errors.New("database error"))
	}
	return nil
}

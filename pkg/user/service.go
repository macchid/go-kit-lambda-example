package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/macchid/learning/userService/pkg/utils"
)

type UserService interface {
	FetchAll() ([]User, error)
	FetchOne(email string) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(email string) (User, error)
	Logger() *log.Logger
}

type userService struct {
	repo   UserData
	logger *log.Logger
}

func NewService(repo UserData, logger *log.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

const (
	errorInvalidUserData    = "invalid user data"
	errorInvalidEmail       = "invalid email"
	errorUserAlreadyExists  = "user already exists"
	errorUserDoesNotExists  = "user does not exist"
	errorUnableToCreateUser = "unable to create user"
	errorUnableToDeleteUser = "unable to delete user"
)

func (svc *userService) FetchAll() ([]User, error) {
	logger, logEnd := utils.LogStart(svc.logger, "UserService::FetchAll")
	defer logEnd(time.Now())

	users, err := svc.repo.FetchAll()
	if err != nil {
		level.Error(logger).Log("msg", "Attempt to fetch all users failed", "err", err)
		return nil, fmt.Errorf("Couldn't get the users list")
	}

	return users, nil
}

func (svc *userService) FetchOne(email string) (User, error) {
	logger, logEnd := utils.LogStart(svc.logger, "UserService::FetchOne")
	defer logEnd(time.Now())

	if !utils.IsEmailValid(email) {
		level.Error(logger).Log("msg", fmt.Sprintf("%s email is invalid", email), "err", errors.New(errorInvalidEmail))
		return User{}, fmt.Errorf("Couldn't find user with email %v", email)
	}

	user, err := svc.repo.FetchOne(email)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Attempt to fetch user with email %v failed", email), "err", err)
		return User{}, fmt.Errorf("Couldn't find user with email %v", email)
	}

	return user, nil
}

func (svc *userService) Create(user User) (User, error) {
	logger, logEnd := utils.LogStart(svc.logger, "UserService::FetchOne")
	defer logEnd(time.Now())

	existent, _ := svc.FetchOne(user.Email)
	if len(existent.Email) != 0 {
		err := errors.New(errorUserAlreadyExists)
		level.Warn(logger).Log("msg", fmt.Sprintf("User with email %v is already registered", user.Email), "err", err)
		return User{}, err
	}

	err := svc.repo.Persist(user)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Couldn't create user %v", user), "err", err)
		return User{}, errors.New(errorUnableToCreateUser)
	}

	return user, nil
}

func (svc *userService) Update(user User) (User, error) {
	logger, logEnd := utils.LogStart(svc.logger, "UserService::Delete")
	defer logEnd(time.Now())

	existent, _ := svc.FetchOne(user.Email)
	if len(existent.Email) == 0 {
		err := errors.New(errorUserDoesNotExists)
		level.Error(logger).Log("msg", fmt.Sprintf("Can't update an inexistent user. Email: %v .", user.Email), "err", err)
		return User{}, err
	}

	if len(user.FirstName) != 0 && user.FirstName != existent.FirstName {
		existent.FirstName = user.FirstName
	}

	if len(user.LastName) != 0 && user.LastName != existent.LastName {
		existent.LastName = user.LastName
	}

	err := svc.repo.Persist(existent)
	if err != nil {
		level.Error(logger).Log("msg", "Couldn't persist the changes in the database.", "err", err)
		return User{}, err
	}

	return existent, nil
}

func (svc *userService) Delete(email string) (User, error) {
	logger, logEnd := utils.LogStart(svc.logger, "UserService::Delete")
	defer logEnd(time.Now())

	existent, _ := svc.FetchOne(email)
	if len(existent.Email) == 0 {
		err := errors.New(errorUserDoesNotExists)
		level.Error(logger).Log("msg", fmt.Sprintf("Won't delete user with email %v. User already deleted.", email), "err", err)
		return User{}, err
	}

	err := svc.repo.Delete(email)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Couldn't erase the client with email %v", email), "err", err)
		return User{}, errors.New(errorUnableToDeleteUser)
	}

	return existent, nil
}

func (svc *userService) Logger() *log.Logger {
	return svc.logger
}

package service

import (
	"../dto"
	"../storage"
	"context"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"log"
	"os"
	"regexp"
)

// третий параметр в функциях - isInternal, для определения типа ошибки в handlers
type UserServiceAPI interface {
	CreateUser(createUserRequest dto.CreateUserRequest) (uuid.UUID, error, bool)
}

type userService struct {
	storage storage.StorageAPI
	ctx context.Context
	log *log.Logger
}

func NewUserServiceAPI(api storage.StorageAPI) UserServiceAPI {
	return &userService{
		storage: api,
		ctx: context.Background(),
		log: log.New(os.Stdout, "USER-SERVICE: ", log.LstdFlags),
	}
}

func (u *userService) CreateUser(createUserRequest dto.CreateUserRequest) (uuid.UUID, error, bool) {
	u.log.Printf("Trying to create user with username: %s", createUserRequest.Username)
	if len(createUserRequest.Username) < 3 {
		return uuid.Nil, xerrors.Errorf("Username must contain at least 3 characters"), false
	}

	var validLogin = regexp.MustCompile("^([a-zA-Z0-9_]+)$")
	f := validLogin.FindStringSubmatch(createUserRequest.Username)
	if f == nil {
		return uuid.Nil, xerrors.Errorf("Username must contain only numbers, latin letters and '_'"), false
	}
	u.log.Printf("Username is valid")

	ok, err := u.storage.GetUserStorage().IsUserExist(createUserRequest.Username)
	if err != nil {
		u.log.Printf("Error while check user on exist in DB, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}
	if ok {
		return uuid.Nil, xerrors.Errorf("User already exist"), false
	}

	tx, err := u.storage.GetTransaction(u.ctx)
	if err != nil {
		u.log.Printf("Error while create transaction, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	id, err := u.storage.GetUserStorage().CreateUser(tx, createUserRequest.Username)
	if err != nil {
		tx.Rollback(u.ctx)
		u.log.Printf("Error while create user in DB, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	err = tx.Commit(u.ctx)
	if err != nil {
		u.log.Printf("Error while commit transaction, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	return id, nil, false
}


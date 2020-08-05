package storage

import (
	"../db"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type UserStorageAPI interface {
	CreateUser(tx pgx.Tx, username string) (uuid.UUID, error)
	IsUserExist(username string) (bool, error)
	CheckExistUsers(ids ...uuid.UUID) (bool, error)
}

type userStorage struct {
	db db.ConnDB
	ctx context.Context
}

func NewUserStorageAPI(connDB db.ConnDB, ctx context.Context) UserStorageAPI {
	return &userStorage{
		db: connDB,
		ctx: ctx,
	}
}

func (u *userStorage) IsUserExist(username string) (bool, error) {
	var result int
	err := u.db.DB.QueryRow(u.ctx, `select count(*) from users where username=$1`, username).Scan(&result)
	if result == 1 {
		return true, err
	}
	return false, nil
}

func (u *userStorage) CreateUser(tx pgx.Tx, username string) (uuid.UUID, error) {
	userID := uuid.Must(uuid.NewUUID())
	if _, err := tx.Exec(u.ctx, `insert into users (id, username) values ($1,$2)`, userID, username); err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (u *userStorage) CheckExistUsers(ids ...uuid.UUID) (bool, error) {
	var result int
	paramsString, userIds := makeParamsFromUUID(ids)
	err := u.db.DB.QueryRow(u.ctx, fmt.Sprintf(`select count(*) from users where id in (%s)`, paramsString), userIds...).Scan(&result)
	if err != nil {
		return false, err
	}

	if result != len(userIds) {
		return false, nil
	}

	return true, nil
}

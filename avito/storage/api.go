package storage

import (
	"../db"
	"context"
	"github.com/jackc/pgx"
)

type StorageAPI interface {
	GetUserStorage() UserStorageAPI
	GetChatStorage() ChatStorageAPI
	GetMessageStorage() MessageStorageAPI
	GetTransaction(ctx context.Context) (pgx.Tx, error)
}

type storageAPI struct {
	userStorage UserStorageAPI
	chatStorage ChatStorageAPI
	messageStorage MessageStorageAPI
	connDB db.ConnDB
}

func (s *storageAPI) GetTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.connDB.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *storageAPI) GetUserStorage() UserStorageAPI {
	return s.userStorage
}

func (s *storageAPI) GetChatStorage() ChatStorageAPI {
	return s.chatStorage
}

func (s *storageAPI) GetMessageStorage() MessageStorageAPI {
	return s.messageStorage
}

func NewStorageAPI(connDB db.ConnDB, ctx context.Context) StorageAPI {
	return &storageAPI{
		userStorage: NewUserStorageAPI(connDB, ctx),
		chatStorage: NewChatStorageAPI(connDB, ctx),
		messageStorage: NewMessageStorageAPI(connDB, ctx),
		connDB: connDB,
	}
}

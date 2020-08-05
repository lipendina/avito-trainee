package service

import (
	"../dto"
	"../storage"
	"context"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"log"
	"os"
)

// третий параметр в функциях - isInternal, для определения типа ошибки в handlers
type ChatServiceAPI interface {
	CreateChat(createChatRequest dto.CreateChatRequest) (uuid.UUID, error, bool)
	GetChatList(chatListRequest dto.ChatListRequest) ([]dto.Chat, error, bool)
}

type chatService struct {
	storage storage.StorageAPI
	ctx context.Context
	log *log.Logger
}

func NewChatServiceAPI(api storage.StorageAPI) ChatServiceAPI {
	return &chatService{
		storage: api,
		ctx: context.Background(),
		log: log.New(os.Stdout, "CHAT-SERVICE: ", log.LstdFlags),
	}
}

func (c *chatService) GetChatList(chatListRequest dto.ChatListRequest) ([]dto.Chat, error, bool) {
	c.log.Printf("Trying to get chats of user %s", chatListRequest)
	ok, err := c.storage.GetUserStorage().CheckExistUsers(chatListRequest.User)
	if err != nil {
		c.log.Printf("Error while check exist users, reason: %+v", err)
		return nil, xerrors.Errorf("System error. Contact support"), true
	}
	if !ok {
		return nil, xerrors.Errorf("User is not exist"), false
	}

	chats, err := c.storage.GetChatStorage().GetChatList(chatListRequest.User)
	if err != nil {
		c.log.Printf("Error while get chats from DB, reason: %+v", err)
		return nil, xerrors.Errorf("System error. Contact support"), true
	}

	return chats, nil, false
}

func (c *chatService) CreateChat(createChatRequest dto.CreateChatRequest) (uuid.UUID, error, bool) {
	c.log.Printf("Trying to create chat: %s", createChatRequest.Name)
	if len(createChatRequest.Name) == 0 {
		return uuid.Nil, xerrors.Errorf("Chat name is empty"), false
	}
	c.log.Printf("Chat name is valid")

	if len(createChatRequest.Users) <= 1 {
		return uuid.Nil, xerrors.Errorf("Not enough users to create chat"), false
	}
	c.log.Printf("Enough users to create chat")

	ok, err := c.storage.GetUserStorage().CheckExistUsers(createChatRequest.Users...)
	if err != nil {
		c.log.Printf("Error while exist users in DB, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}
	if !ok {
		return uuid.Nil, xerrors.Errorf("One or more users are not exist"), false
	}

	tx, err := c.storage.GetTransaction(c.ctx)
	if err != nil {
		c.log.Printf("Error while create transaction, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	chatID, err := c.storage.GetChatStorage().CreateChat(tx, createChatRequest.Name)
	if err != nil {
		c.log.Printf("Error while create chat in DB, reason: %+v", err)
		tx.Rollback(c.ctx)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	err = c.storage.GetChatStorage().CreateRecordChatsUsers(tx, chatID, createChatRequest.Users...)
	if err != nil {
		c.log.Printf("Error while create record in chats_users, reason: %+v", err)
		tx.Rollback(c.ctx)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	err = tx.Commit(c.ctx)
	if err != nil {
		c.log.Printf("Error while commit transaction, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	return chatID, nil, false
}

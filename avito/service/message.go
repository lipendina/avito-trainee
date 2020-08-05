package service

import (
	"../dto"
	"../storage"
	"context"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"log"
	"os"
	"strings"
)

// третий параметр в функциях - isInternal, для определения типа ошибки в handlers
type MessageServiceAPI interface {
	SendMessage(sendMessageRequest dto.SendMessageRequest) (uuid.UUID, error, bool)
	GetMessageList(getMessageList dto.MessageListRequest) ([]dto.Message, error, bool)
}

type messageService struct {
	storage storage.StorageAPI
	ctx context.Context
	log *log.Logger
}

func NewMessageServiceAPI(api storage.StorageAPI) MessageServiceAPI {
	return &messageService{
		storage: api,
		ctx: context.Background(),
		log: log.New(os.Stdout, "MESSAGE-SERVICE: ", log.LstdFlags),
	}
}

func (m *messageService) SendMessage(sendMessageRequest dto.SendMessageRequest) (uuid.UUID, error, bool) {
	m.log.Printf("Trying to send message: %s", sendMessageRequest)
	// constraint по user_id и chat_id гарантируют, что сущности существуют
	ok, err := m.storage.GetMessageStorage().CheckExistUserChats(sendMessageRequest.Author, sendMessageRequest.Chat)
	if err != nil {
		m.log.Printf("Error while check exist user in chat, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}
	if !ok {
		return uuid.Nil, xerrors.Errorf("User doesn't consist in chat"), false
	}
	m.log.Printf("Author of message exist in chat")

	if len(strings.TrimSpace(sendMessageRequest.Text)) == 0 {
		return uuid.Nil, xerrors.Errorf("Empty message"), false
	}

	tx, err := m.storage.GetTransaction(m.ctx)
	if err != nil {
		m.log.Printf("Error while create transaction, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	messageID, err := m.storage.GetMessageStorage().CreateMessage(tx, sendMessageRequest.Author, sendMessageRequest.Chat, sendMessageRequest.Text)
	if err != nil {
		m.log.Printf("Error while create message, reason: %+v", err)
		tx.Rollback(m.ctx)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	err = tx.Commit(m.ctx)
	if err != nil {
		m.log.Printf("Error while commit transaction, reason: %+v", err)
		return uuid.Nil, xerrors.Errorf("System error. Contact support"), true
	}

	return messageID, nil, false
}

func (m *messageService) GetMessageList(getMessageList dto.MessageListRequest) ([]dto.Message, error, bool) {
	m.log.Printf("Trying to get messages in chat: %s", getMessageList)
	ok, err := m.storage.GetChatStorage().CheckExistChat(getMessageList.Chat)
	if err != nil {
		m.log.Printf("Error while check exist chat, reason: %+v", err)
		return nil, xerrors.Errorf("System error. Contact support"), true
	}
	if !ok {
		return nil, xerrors.Errorf("Chat doesn't exist"), false
	}
	m.log.Printf("Chat is exist")

	messages, err := m.storage.GetMessageStorage().GetMessageList(getMessageList.Chat)
	if err != nil {
		m.log.Printf("Error while get message list, reason: %+v", err)
		return nil, xerrors.Errorf("System error. Contact support"), true
	}

	return messages, nil, false
}

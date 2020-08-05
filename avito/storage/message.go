package storage

import (
	"../dto"
	"../db"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type MessageStorageAPI interface {
	CreateMessage(tx pgx.Tx, author uuid.UUID, chat uuid.UUID, text string) (uuid.UUID, error)
	CheckExistUserChats(author uuid.UUID, chat uuid.UUID) (bool, error)
	GetMessageList(chat uuid.UUID) ([]dto.Message, error)
}

type messageStorage struct {
	db db.ConnDB
	ctx context.Context
}

func NewMessageStorageAPI(connDB db.ConnDB, ctx context.Context) MessageStorageAPI {
	return &messageStorage{
		db: connDB,
		ctx: ctx,
	}
}

func (m *messageStorage) GetMessageList(chat uuid.UUID) ([]dto.Message, error) {
	rows, err := m.db.DB.Query(context.Background(), `select id, chat, author, text, extract(epoch from created_at) as created_at from messages where chat=$1 order by created_at asc`, chat)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	messages := make([]dto.Message, 0)
	for rows.Next() {
		var message dto.Message
		err := rows.Scan(&message.ID, &message.Chat, &message.Author, &message.Text, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil

}

func (m *messageStorage) CreateMessage(tx pgx.Tx, author uuid.UUID, chat uuid.UUID, text string) (uuid.UUID, error) {
	messageID := uuid.Must(uuid.NewUUID())
	if _, err := tx.Exec(context.Background(), `insert into messages (id, chat, author, text) values ($1, $2, $3, $4)`,
		messageID, chat, author, text); err != nil {
		return uuid.Nil, err
	}

	return messageID, nil
}

func (m *messageStorage) CheckExistUserChats(author uuid.UUID, chat uuid.UUID) (bool, error) {
	var result int
	err := m.db.DB.QueryRow(context.Background(), `select count(*) from chats_users where user_id=$1 and chat_id=$2`, author, chat).Scan(&result)
	if err != nil {
		return false, err
	}

	if result != 1 {
		return false, nil
	}

	return true, nil
}

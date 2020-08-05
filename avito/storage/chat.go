package storage

import (
	"../db"
	"../dto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"strings"
)

type ChatStorageAPI interface {
	CreateChat(tx pgx.Tx, chatname string) (uuid.UUID, error)
	CreateRecordChatsUsers(tx pgx.Tx, chatID uuid.UUID, users ...uuid.UUID) error
	GetChatList(userId uuid.UUID) ([]dto.Chat, error)
	CheckExistChat(chat uuid.UUID) (bool, error)
}

type chatStorage struct {
	db db.ConnDB
	ctx context.Context
}

func NewChatStorageAPI(connDB db.ConnDB, ctx context.Context) ChatStorageAPI {
	return &chatStorage{
		db: connDB,
		ctx: ctx,
	}
}

func (c *chatStorage) CreateChat(tx pgx.Tx, chatname string) (uuid.UUID, error) {
	chatID := uuid.Must(uuid.NewUUID())
	if _, err := tx.Exec(c.ctx, `insert into chats (id, name) values ($1, $2)`, chatID, chatname); err != nil {
		return uuid.Nil, err
	}

	return chatID, nil
}

func (c *chatStorage) CreateRecordChatsUsers(tx pgx.Tx, chatID uuid.UUID, users ...uuid.UUID) error {
	valueString := make([]string, 0, len(users))
	valueArgs := make([]interface{}, 0, len(users) * 3)
	for idx, user := range users {
		recordID := uuid.Must(uuid.NewUUID())
		valueString = append(valueString, fmt.Sprintf("($%d,$%d,$%d)", idx*3+1, idx*3+2, idx*3+3))
		valueArgs = append(valueArgs, recordID)
		valueArgs = append(valueArgs, user)
		valueArgs = append(valueArgs, chatID)
	}

	_, err := tx.Exec(c.ctx, fmt.Sprintf(`insert into chats_users (id, user_id, chat_id) values %s`, strings.Join(valueString, ",")), valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (c *chatStorage) GetChatList(userId uuid.UUID) ([]dto.Chat, error) {
	rows, err := c.db.DB.Query(c.ctx, `select t1.chat_id as chat_id, t1.name as name, extract(epoch from coalesce(t2.created_at, t1.created_at)) as created_at 
from (select u.chat_id, c.name, c.created_at from chats_users u join chats c 
	on u.chat_id = c.id where user_id=$1) t1 
left join (select chat, created_at from messages order by created_at desc limit 1) t2 
	on t1.chat_id = t2.chat 
order by created_at desc`, userId)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	chats := make([]dto.Chat, 0)
	chatIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		var chat dto.Chat
		err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}

		chatIDs = append(chatIDs, chat.ID)

		chats = append(chats, chat)
	}

	usersByChatID := make(map[uuid.UUID][]uuid.UUID)
	paramsString, parsedIDs := makeParamsFromUUID(chatIDs)
	rows, err = c.db.DB.Query(c.ctx, fmt.Sprintf(`select chat_id, user_id from chats_users where chat_id in (%s)`, paramsString), parsedIDs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var chatID, userID uuid.UUID
		err := rows.Scan(&chatID, &userID)
		if err != nil {
			return nil, err
		}
		usersByChatID[chatID] = append(usersByChatID[chatID], userID)
	}

	for i := range chats {
		chats[i].Users = usersByChatID[chats[i].ID]
	}

	return chats, nil
}

func (c *chatStorage) CheckExistChat(chat uuid.UUID) (bool, error) {
	var result int
	err := c.db.DB.QueryRow(context.Background(), `select count(*) from chats where id=$1`, chat).Scan(&result)
	if err != nil {
		return false, err
	}

	return true, nil
}
package dto

import (
	"fmt"
	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	Name      string
	Users     []uuid.UUID
	CreatedAt float64
}

type CreateChatRequest struct {
	Name  string      `json: "name"`
	Users []uuid.UUID `json: "users"`
}

func (r CreateChatRequest) String() string {
	return fmt.Sprintf("{chatname: %s}", r.Name)
}

type CreateChatResponse struct {
	ID uuid.UUID `json: "id"`
}

func (r CreateChatResponse) String() string {
	return fmt.Sprintf("{chatID: %s}", r.ID)
}

type ChatListRequest struct {
	User uuid.UUID `json: "user"`
}

func (r ChatListRequest) String() string {
	return fmt.Sprintf("{user: %s}", r.User)
}

type ChatListResponse struct {
	ChatList []Chat
}

func (r ChatListResponse) String() string {
	return fmt.Sprintf("{chats: %s}", r.ChatList)
}
package dto

import (
	"fmt"
	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	Chat      uuid.UUID
	Author    uuid.UUID
	Text      string
	CreatedAt float64
}

type SendMessageRequest struct {
	Chat   uuid.UUID `json: "chat"`
	Author uuid.UUID `json: "author"`
	Text   string    `json: "text"`
}

func (r SendMessageRequest) String() string {
	return fmt.Sprintf("{authorID: %s, chatID: %s, text: %s}", r.Author, r.Chat, r.Text)
}

type SendMessageResponse struct {
	ID uuid.UUID `json: "id"`
}

func (r SendMessageResponse) String() string {
	return fmt.Sprintf("{messageID: %s}", r.ID)
}

type MessageListRequest struct {
	Chat uuid.UUID `json: "chat"`
}

func (r MessageListRequest) String() string {
	return fmt.Sprintf("{chatID: %s}", r.Chat)
}

type MessageListResponse struct {
	MessageList []Message
}

func (r MessageListResponse) String() string {
	return fmt.Sprintf("{messages: %s}", r.MessageList)
}

package dto

import (
	"fmt"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string `json: "username"`
}

func (r CreateUserRequest) String() string {
	return fmt.Sprintf("{username: %s}", r.Username)
}

type CreateUserResponse struct {
	ID uuid.UUID `json: "id"`
}

func (r CreateUserResponse) String() string {
	return fmt.Sprintf("{userID: %s}", r.ID)
}
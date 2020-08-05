package service

import (
	"../storage"
)

type ServiceAPI interface {
	GetUserService() UserServiceAPI
	GetChatService() ChatServiceAPI
	GetMessageService() MessageServiceAPI
}

type serviceAPI struct {
	userServiceAPI UserServiceAPI
	chatServiceAPI ChatServiceAPI
	messageServiceAPI MessageServiceAPI
}

func NewServiceAPI(api storage.StorageAPI) ServiceAPI {
	return &serviceAPI{
		userServiceAPI: NewUserServiceAPI(api),
		chatServiceAPI: NewChatServiceAPI(api),
		messageServiceAPI: NewMessageServiceAPI(api),
	}
}

func (s *serviceAPI) GetUserService() UserServiceAPI {
	return s.userServiceAPI
}

func (s *serviceAPI) GetChatService() ChatServiceAPI {
	return s.chatServiceAPI
}

func (s *serviceAPI) GetMessageService() MessageServiceAPI {
	return s.messageServiceAPI
}

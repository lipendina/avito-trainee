package handlers

import (
	"../dto"
	"../service"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Handlers interface {
	AddNewUserHandler(w http.ResponseWriter, r *http.Request)
	CreateChatHandler(w http.ResponseWriter, r *http.Request)
	SendMessageHandler(w http.ResponseWriter, r *http.Request)

	GetChatListHandler(w http.ResponseWriter, r *http.Request)
	GetMessageListHandler(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	service service.ServiceAPI
	log *log.Logger
}

func NewHandlers(api service.ServiceAPI) Handlers {
	return &handlers{
		service: api,
		log: log.New(os.Stdout, "CONTROLLER: ", log.LstdFlags),
	}
}

func (h *handlers) AddNewUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createUserRequest dto.CreateUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&createUserRequest)

	if err != nil {
		h.log.Printf("Error while parse createUserRequest, reason: %v", err)
		response := &dto.ErrorResponse{Message: "Cannot parse request body"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}
	h.log.Printf("Received createUserRequest: %s", createUserRequest)

	userID, err, isInternal := h.service.GetUserService().CreateUser(createUserRequest)
	if err != nil {
		h.log.Printf("Error while createUser, reason: %v", err)
		response := &dto.ErrorResponse{Message: err.Error()}
		sendResponse(getErrorStatus(isInternal), response, w)
		return
	}

	response := &dto.CreateUserResponse{ID: userID}
	h.log.Printf("Send response: %v", response)
	sendResponse(http.StatusOK, response, w)
}

func (h *handlers) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createChatRequest dto.CreateChatRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&createChatRequest)
	if err != nil {
		h.log.Printf("Error while parse createChatRequest, reason: %v", err)
		response := &dto.ErrorResponse{Message: "Cannot parse request"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}
	h.log.Printf("Received createChatRequest: %s", createChatRequest)

	chatID, err, isInternal := h.service.GetChatService().CreateChat(createChatRequest)
	if err != nil {
		h.log.Printf("Error while createChat, reason: %v", err)
		response := &dto.ErrorResponse{Message: err.Error()}
		sendResponse(getErrorStatus(isInternal), response, w)
		return
	}

	response := &dto.CreateChatResponse{ID: chatID}
	h.log.Printf("Send response: %v", response)
	sendResponse(http.StatusOK, response, w)
}

func (h *handlers) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var sendMessageRequest dto.SendMessageRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&sendMessageRequest)

	if err != nil {
		h.log.Printf("Error while parse sendMessageRequest, reason: %v", err)
		response := &dto.ErrorResponse{Message: "Cannot parse request"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}
	h.log.Printf("Received sendMessageRequest: %s", sendMessageRequest)

	messageID, err, isInternal := h.service.GetMessageService().SendMessage(sendMessageRequest)
	if err != nil {
		h.log.Printf("Error while send message, reason: %v", err)
		response := &dto.ErrorResponse{Message: err.Error()}
		sendResponse(getErrorStatus(isInternal), response, w)
		return
	}

	response := &dto.SendMessageResponse{ID: messageID}
	h.log.Printf("Send request: %v", sendMessageRequest)
	sendResponse(http.StatusOK, response, w)
}

func (h *handlers) GetChatListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var chatListRequest dto.ChatListRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&chatListRequest)

	if err != nil {
		h.log.Printf("Error while parse chatListRequest, reason, %v", err)
		response := &dto.ErrorResponse{Message: "Cannot parse request"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}
	h.log.Printf("Received chatListRequest: %s", chatListRequest)

	chats, err, isInternal := h.service.GetChatService().GetChatList(chatListRequest)
	if err != nil {
		h.log.Printf("Error while getChatList, reason: %v", err)
		response := &dto.ErrorResponse{Message: err.Error()}
		sendResponse(getErrorStatus(isInternal), response, w)
		return
	}

	response := &dto.ChatListResponse{ChatList: chats}
	h.log.Printf("Send response: %v", response)
	sendResponse(http.StatusOK, response, w)
}

func (h *handlers) GetMessageListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var messageListRequest dto.MessageListRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&messageListRequest)

	if err != nil {
		h.log.Printf("Error while parse getMessageListRequest, reason: %v", err)
		response := &dto.ErrorResponse{Message: "Cannot parse request"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}
	h.log.Printf("Received messageListRequest: %s", messageListRequest)

	messages, err, isInternal := h.service.GetMessageService().GetMessageList(messageListRequest)
	if err != nil {
		h.log.Printf("Error while getMessageList, reason: %v", err)
		response := &dto.ErrorResponse {Message: err.Error()}
		sendResponse(getErrorStatus(isInternal), response, w)
		return
	}

	response := &dto.MessageListResponse{MessageList: messages}
	h.log.Printf("Send response: %v", response)
	sendResponse(http.StatusOK, response, w)
}

func sendResponse(httpStatus int, response interface{}, w http.ResponseWriter) {
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}


package main

import (
	"./config"
	"./db"
	"./handlers"
	"./service"
	"./storage"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	applicationConfig, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Cannot parse config: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pgConn := db.NewConnectToPG(&applicationConfig.DB, ctx)

	storageAPI := storage.NewStorageAPI(pgConn, ctx)
	serviceAPI := service.NewServiceAPI(storageAPI)

	a := handlers.NewHandlers(serviceAPI)

	r := mux.NewRouter()
	// добавление нового пользователя
	r.HandleFunc("/users/add", a.AddNewUserHandler).Methods("POST")
	// создание чата между пользователями
	r.HandleFunc("/chats/add", a.CreateChatHandler).Methods("POST")
	// отправление сообщения от лица пользователя
	r.HandleFunc("/messages/add", a.SendMessageHandler).Methods("POST")
	// получение списка чатов конкретного пользователя
	r.HandleFunc("/chats/get", a.GetChatListHandler).Methods("POST")
	// получение списка сообщений конкретного чата
	r.HandleFunc("/messages/get", a.GetMessageListHandler).Methods("POST")
	http.Handle("/", r)

	fmt.Println("Server is listening...")
	http.ListenAndServe(fmt.Sprintf(":%d", applicationConfig.HTTPPort), nil)
}

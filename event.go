package main

import "log"

type (
	NewNotificationChanParam struct {
		UserID  string `json:"user_id"`
		Message string `json:"message"`
	}

	ClientNotificationChannel    chan NewNotificationChanParam
	ClientNotificationChannelSet map[ClientNotificationChannel]bool

	ChannelLifecycleEvent struct {
		UserID  string
		Channel ClientNotificationChannel
	}

	EventManager struct {
		Message                 ClientNotificationChannel
		NewClients              chan ChannelLifecycleEvent
		ClosedClients           chan ChannelLifecycleEvent
		UserNotificationClients map[string]ClientNotificationChannelSet
	}
)

var eventManager *EventManager

func NewEventManager() *EventManager {
	eventManager = &EventManager{
		Message:                 make(ClientNotificationChannel),
		NewClients:              make(chan ChannelLifecycleEvent),
		ClosedClients:           make(chan ChannelLifecycleEvent),
		UserNotificationClients: make(map[string]ClientNotificationChannelSet),
	}

	go eventManager.listen()

	return eventManager
}

func GetEventManager() *EventManager {
	if eventManager == nil {
		eventManager = NewEventManager()
	}

	return eventManager
}

func (stream *EventManager) listen() {
	for {
		select {
		case client := <-stream.NewClients:
			selectedNotificationClient, ok := stream.UserNotificationClients[client.UserID]
			if !ok {
				selectedNotificationClient = make(ClientNotificationChannelSet)
			}
			selectedNotificationClient[client.Channel] = true

			stream.UserNotificationClients[client.UserID] = selectedNotificationClient

			log.Printf("Channel for user %s added. %d registered clients\n", client.UserID, len(selectedNotificationClient))
		case client := <-stream.ClosedClients:
			close(client.Channel)
			delete(stream.UserNotificationClients[client.UserID], client.Channel)
			log.Printf("Channel for user %s closed. %d registered clients\n", client.UserID, len(stream.UserNotificationClients[client.UserID]))
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.UserNotificationClients[eventMsg.UserID] {
				clientMessageChan <- eventMsg
			}
		}
	}
}

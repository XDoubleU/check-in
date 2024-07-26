package services

import (
	"check-in/api/internal/models"
	"sync"

	"nhooyr.io/websocket"
)

type WebSocketService struct {
	subscribers map[*websocket.Conn]models.Subscriber
}

func (service WebSocketService) GetAllUpdateEvents(
	conn *websocket.Conn,
) []models.LocationUpdateEvent {
	for {
		if service.subscribers[conn].BufferMu.TryLock() {
			break
		}
	}
	defer service.subscribers[conn].BufferMu.Unlock()

	var result []models.LocationUpdateEvent

	for key, event := range service.subscribers[conn].Buffer {
		result = append(result, event)
		delete(service.subscribers[conn].Buffer, key)
	}

	return result
}

func (service WebSocketService) GetByNormalizedName(
	conn *websocket.Conn,
) models.LocationUpdateEvent {
	for {
		if service.subscribers[conn].BufferMu.TryLock() {
			break
		}
	}
	defer service.subscribers[conn].BufferMu.Unlock()

	name := service.subscribers[conn].NormalizedName

	result := service.subscribers[conn].Buffer[name]
	delete(service.subscribers[conn].Buffer, name)

	return result
}

func (service WebSocketService) AddUpdateEvent(location models.Location) {
	locationUpdateEvent := models.LocationUpdateEvent{
		NormalizedName:     location.NormalizedName,
		Available:          location.Available,
		Capacity:           location.Capacity,
		YesterdayFullAt:    location.YesterdayFullAt,
		AvailableYesterday: location.AvailableYesterday,
		CapacityYesterday:  location.CapacityYesterday,
	}

	for _, subscriber := range service.subscribers {
		if !(subscriber.Subject == "all-locations" ||
			(subscriber.Subject == "single-location" &&
				subscriber.NormalizedName == locationUpdateEvent.NormalizedName)) {
			continue
		}

		for {
			if subscriber.BufferMu.TryLock() {
				break
			}
		}
		defer subscriber.BufferMu.Unlock()

		subscriber.Buffer[locationUpdateEvent.NormalizedName] = locationUpdateEvent
	}
}

func (service WebSocketService) AddSubscriber(
	conn *websocket.Conn,
	subject models.WebSocketSubject,
	normalizedName string,
) {
	var mu sync.Mutex

	service.subscribers[conn] = models.Subscriber{
		Subject:        subject,
		NormalizedName: normalizedName,
		Buffer:         make(map[string]models.LocationUpdateEvent),
		BufferMu:       &mu,
	}
}

func (service WebSocketService) RemoveSubscriber(conn *websocket.Conn) {
	delete(service.subscribers, conn)
}

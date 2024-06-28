package repositories

import (
	"sync"

	"nhooyr.io/websocket"

	"check-in/api/internal/models"
)

type WebSocketRepository struct {
	subscribers map[*websocket.Conn]models.Subscriber
}

func (repo WebSocketRepository) GetAllUpdateEvents(
	conn *websocket.Conn,
) []models.LocationUpdateEvent {
	for {
		if repo.subscribers[conn].BufferMu.TryLock() {
			break
		}
	}
	defer repo.subscribers[conn].BufferMu.Unlock()

	var result []models.LocationUpdateEvent

	for key, event := range repo.subscribers[conn].Buffer {
		result = append(result, event)
		delete(repo.subscribers[conn].Buffer, key)
	}

	return result
}

func (repo WebSocketRepository) GetByNormalizedName(
	conn *websocket.Conn,
) models.LocationUpdateEvent {
	for {
		if repo.subscribers[conn].BufferMu.TryLock() {
			break
		}
	}
	defer repo.subscribers[conn].BufferMu.Unlock()

	name := repo.subscribers[conn].NormalizedName

	result := repo.subscribers[conn].Buffer[name]
	delete(repo.subscribers[conn].Buffer, name)

	return result
}

func (repo WebSocketRepository) AddUpdateEvent(location models.Location) {
	locationUpdateEvent := models.LocationUpdateEvent{
		NormalizedName:     location.NormalizedName,
		Available:          location.Available,
		Capacity:           location.Capacity,
		YesterdayFullAt:    location.YesterdayFullAt,
		AvailableYesterday: location.AvailableYesterday,
		CapacityYesterday:  location.CapacityYesterday,
	}

	for _, subscriber := range repo.subscribers {
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

func (repo WebSocketRepository) AddSubscriber(
	conn *websocket.Conn,
	subject models.WebSocketSubject,
	normalizedName string,
) {
	var mu sync.Mutex

	repo.subscribers[conn] = models.Subscriber{
		Subject:        subject,
		NormalizedName: normalizedName,
		Buffer:         make(map[string]models.LocationUpdateEvent),
		BufferMu:       &mu,
	}
}

func (repo WebSocketRepository) RemoveSubscriber(conn *websocket.Conn) {
	delete(repo.subscribers, conn)
}

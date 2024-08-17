package services

import (
	"context"
	"net/http"

	wstools "github.com/xdoubleu/essentia/pkg/communication/ws"
	errortools "github.com/xdoubleu/essentia/pkg/errors"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

type GetAllLocationStatesFunc = func(ctx context.Context) ([]dtos.LocationStateDto, error)
type GetStateFunc = func(ctx context.Context) (*models.State, error)

type WebSocketService struct {
	handler           *wstools.WebSocketHandler[dtos.SubscribeMessageDto]
	stateTopic        *wstools.Topic
	allLocationsTopic *wstools.Topic
	locationTopics    map[string]*wstools.Topic
}

func NewWebSocketService(
	allowedOrigin string,
) *WebSocketService {
	service := WebSocketService{
		handler:           nil,
		allLocationsTopic: nil,
		locationTopics:    make(map[string]*wstools.Topic),
	}

	handler := wstools.CreateWebSocketHandler[dtos.SubscribeMessageDto](
		1,
		100, //nolint:mnd //no magic number
		[]string{allowedOrigin},
	)
	service.handler = &handler

	return &service
}

func (service WebSocketService) Handler() http.HandlerFunc {
	return service.handler.Handler()
}

func (service *WebSocketService) SetStateTopic(getState GetStateFunc) error {
	topic, err := service.handler.AddTopic(
		string(dtos.State),
		func(ctx context.Context, _ *wstools.Topic) (any, error) { return getState(ctx) },
	)
	if err != nil {
		return err
	}

	service.stateTopic = topic
	return nil
}

func (service *WebSocketService) SetAllLocationsTopic(getAllLocationStates GetAllLocationStatesFunc) error {
	topic, err := service.handler.AddTopic(
		"*",
		func(ctx context.Context, _ *wstools.Topic) (any, error) { return getAllLocationStates(ctx) },
	)
	if err != nil {
		return err
	}

	service.allLocationsTopic = topic
	return nil
}

func (service WebSocketService) AddLocation(location *models.Location) error {
	topic, err := service.handler.AddTopic(location.NormalizedName, nil)
	if err != nil {
		return err
	}

	service.locationTopics[location.ID] = topic
	return nil
}

func (service WebSocketService) UpdateLocation(location *models.Location) error {
	topic, ok := service.locationTopics[location.ID]
	if !ok {
		return errortools.NewNotFoundError("location", location.ID, "id")
	}

	newTopic, err := service.handler.UpdateTopicName(topic, location.NormalizedName)
	if err != nil {
		return err
	}

	delete(service.locationTopics, location.ID)
	service.locationTopics[location.ID] = newTopic
	return nil
}

func (service WebSocketService) DeleteLocation(location *models.Location) error {
	topic, ok := service.locationTopics[location.ID]
	if !ok {
		return errortools.NewNotFoundError("location", location.ID, "id")
	}

	err := service.handler.RemoveTopic(topic)
	if err != nil {
		return err
	}

	delete(service.locationTopics, topic.Name)
	return nil
}

func (service WebSocketService) NewAppState(state models.State) {
	service.stateTopic.EnqueueEvent(state)
}

func (service WebSocketService) NewLocationState(location models.Location) {
	locationState := dtos.NewLocationStateDto(location)

	service.allLocationsTopic.EnqueueEvent(locationState)
	service.locationTopics[location.ID].EnqueueEvent(locationState)
}

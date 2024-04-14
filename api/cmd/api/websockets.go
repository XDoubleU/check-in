package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"check-in/api/internal/config"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/validator"
)

func (app *application) websocketsRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/", app.webSocketHandler)
}

// @Summary	WebSocket for receiving update events
// @Tags		websocket
// @Param		subscribeMessageDto	body		SubscribeMessageDto	true	"SubscribeMessageDto"
// @Success	200					{object}	LocationUpdateEvent
// @Router		/ws [get].
func (app *application) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	url := app.config.WebURL
	if strings.Contains(url, "://") {
		url = strings.Split(app.config.WebURL, "://")[1]
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{url},
	})
	if err != nil {
		return
	}
	defer app.unsubAndClose(conn)

	var msg dtos.SubscribeMessageDto
	err = wsjson.Read(r.Context(), conn, &msg)
	if err != nil {
		app.handleWsError(r.Context(), conn, err)
		return
	}

	v := validator.New()

	if dtos.ValidateSubscribeMessageDto(v, msg); !v.Valid() {
		app.handleWsError(r.Context(), conn, errors.New(v.Errors["normalizedName"]))
		return
	}

	switch {
	case msg.Subject == models.AllLocations:
		allLocationsHandler(r.Context(), app, conn, msg)

	case msg.Subject == models.SingleLocation:
		singleLocationHandler(r.Context(), app, conn, msg)

	default:
		return
	}
}

func allLocationsHandler(
	ctx context.Context,
	app *application,
	conn *websocket.Conn,
	msg dtos.SubscribeMessageDto,
) {
	app.services.WebSockets.AddSubscriber(conn, msg.Subject, msg.NormalizedName)

	locationUpdateEventDtos, _ := app.getAllCurrentLocationStates(ctx)

	err := wsjson.Write(ctx, conn, locationUpdateEventDtos)
	if err != nil {
		app.handleWsError(ctx, conn, err)
		return
	}

	for {
		updateEvents := app.services.WebSockets.GetAllUpdateEvents(conn)
		if len(updateEvents) > 0 {
			err = wsjson.Write(ctx, conn, updateEvents)
			if err != nil {
				app.handleWsError(ctx, conn, err)
				return
			}
		}

		if app.config.Env != config.TestEnv {
			time.Sleep(30 * time.Second) //nolint:gomnd //no magic number
		}
	}
}

func singleLocationHandler(
	ctx context.Context,
	app *application,
	conn *websocket.Conn,
	msg dtos.SubscribeMessageDto,
) {
	app.services.WebSockets.AddSubscriber(conn, msg.Subject, msg.NormalizedName)

	for {
		updateEvent := app.services.WebSockets.GetByNormalizedName(conn)
		if updateEvent.NormalizedName == msg.NormalizedName {
			err := wsjson.Write(ctx, conn, updateEvent)
			if err != nil {
				app.handleWsError(ctx, conn, err)
				return
			}
		}

		if app.config.Env != config.TestEnv {
			time.Sleep(time.Second)
		}
	}
}

func (app *application) handleWsError(
	ctx context.Context,
	conn *websocket.Conn,
	err error,
) {
	if websocket.CloseStatus(err) != websocket.StatusNormalClosure &&
		websocket.CloseStatus(err) != websocket.StatusGoingAway {
		app.unsubAndClose(conn)
		err = wsjson.Write(ctx, conn, err)
		if err != nil {
			app.logError(err)
		}
		return
	}
}

func (app *application) unsubAndClose(conn *websocket.Conn) {
	app.services.WebSockets.RemoveSubscriber(conn)
	conn.Close(websocket.StatusInternalError, "")
}

func (app *application) getAllCurrentLocationStates(
	ctx context.Context,
) ([]models.LocationUpdateEvent, error) {
	locations, err := app.services.Locations.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	locationUpdateEvents := []models.LocationUpdateEvent{}

	for _, location := range locations {
		locationUpdateEvent := models.LocationUpdateEvent{
			NormalizedName:     location.NormalizedName,
			Available:          location.Available,
			Capacity:           location.Capacity,
			YesterdayFullAt:    location.YesterdayFullAt,
			AvailableYesterday: location.AvailableYesterday,
			CapacityYesterday:  location.CapacityYesterday,
		}

		locationUpdateEvents = append(
			locationUpdateEvents,
			locationUpdateEvent,
		)
	}

	return locationUpdateEvents, nil
}

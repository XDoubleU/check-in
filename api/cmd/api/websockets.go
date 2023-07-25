package main

import (
	"context"
	"net/http"
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
// @Success	200	{object}	LocationUpdateEvent
// @Router		/ws [get].
func (app *application) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer unsubAndClose(conn, app)

	var msg dtos.SubscribeMessageDto
	err = wsjson.Read(r.Context(), conn, &msg)
	if err != nil {
		handleError(err)
		return
	}

	v := validator.New()

	if dtos.ValidateSubscribeMessageDto(v, msg); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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
		handleError(err)
		return
	}

	for {
		updateEvents := app.services.WebSockets.GetAllUpdateEvents(conn)
		if len(updateEvents) > 0 {
			err = wsjson.Write(ctx, conn, updateEvents)
			if err != nil {
				handleError(err)
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
				handleError(err)
				return
			}
		}

		if app.config.Env != config.TestEnv {
			time.Sleep(time.Second)
		}
	}
}

func handleError(err error) {
	if websocket.CloseStatus(err) != websocket.StatusNormalClosure &&
		websocket.CloseStatus(err) != websocket.StatusGoingAway {
		return
	}
}

func unsubAndClose(conn *websocket.Conn, app *application) {
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
			NormalizedName:  location.NormalizedName,
			Available:       location.Available,
			Capacity:        location.Capacity,
			YesterdayFullAt: location.YesterdayFullAt,
		}

		locationUpdateEvents = append(
			locationUpdateEvents,
			locationUpdateEvent,
		)
	}

	return locationUpdateEvents, nil
}

package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/julienschmidt/httprouter"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"check-in/api/internal/config"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
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
		http_tools.WSUpgradeErrorResponse(w, r, err)
		return
	}
	defer func() {
		//todo
		app.services.WebSockets.RemoveSubscriber(conn)
		conn.Close(websocket.StatusInternalError, "")
	}()

	var msg dtos.SubscribeMessageDto
	err = wsjson.Read(r.Context(), conn, &msg)
	if err != nil {
		http_tools.WSErrorResponse(w, r, conn, app.services.WebSockets.RemoveSubscriber, err)
		return
	}

	if v := msg.Validate(); !v.Valid() {
		http_tools.FailedValidationResponse(w, r, v.Errors)
		return
	}

	switch {
	case msg.Subject == models.AllLocations:
		allLocationsHandler(w, r, app, conn, msg)

	case msg.Subject == models.SingleLocation:
		singleLocationHandler(w, r, app, conn, msg)

	default:
		return
	}
}

func allLocationsHandler(
	w http.ResponseWriter,
	r *http.Request,
	app *application,
	conn *websocket.Conn,
	msg dtos.SubscribeMessageDto,
) {
	app.services.WebSockets.AddSubscriber(conn, msg.Subject, msg.NormalizedName)

	locationUpdateEventDtos, _ := app.getAllCurrentLocationStates(r.Context())

	err := wsjson.Write(r.Context(), conn, locationUpdateEventDtos)
	if err != nil {
		http_tools.WSErrorResponse(w, r, conn, app.services.WebSockets.RemoveSubscriber, err)
		return
	}

	for {
		updateEvents := app.services.WebSockets.GetAllUpdateEvents(conn)
		if len(updateEvents) > 0 {
			err = wsjson.Write(r.Context(), conn, updateEvents)
			if err != nil {
				http_tools.WSErrorResponse(w, r, conn, app.services.WebSockets.RemoveSubscriber, err)
				return
			}
		}

		if app.config.Env != config.TestEnv {
			time.Sleep(30 * time.Second) //nolint:gomnd //no magic number
		}
	}
}

func singleLocationHandler(
	w http.ResponseWriter,
	r *http.Request,
	app *application,
	conn *websocket.Conn,
	msg dtos.SubscribeMessageDto,
) {
	app.services.WebSockets.AddSubscriber(conn, msg.Subject, msg.NormalizedName)

	for {
		updateEvent := app.services.WebSockets.GetByNormalizedName(conn)
		if updateEvent.NormalizedName == msg.NormalizedName {
			err := wsjson.Write(r.Context(), conn, updateEvent)
			if err != nil {
				http_tools.WSErrorResponse(w, r, conn, app.services.WebSockets.RemoveSubscriber, err)
				return
			}
		}

		if app.config.Env != config.TestEnv {
			time.Sleep(time.Second)
		}
	}
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

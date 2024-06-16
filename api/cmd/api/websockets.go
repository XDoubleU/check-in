package main

import (
	"context"
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/julienschmidt/httprouter"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"check-in/api/internal/config"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

// @Summary	WebSocket for receiving update events
// @Tags		websocket
// @Param		subscribeMessageDto	body		SubscribeMessageDto	true	"SubscribeMessageDto"
// @Success	200					{object}	LocationUpdateEvent
// @Router		/ws [get].
func (app *application) websocketsRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/", app.getWebSocketHandler())
}

func (app *application) getWebSocketHandler() http.HandlerFunc {
	wsh := http_tools.CreateWebsocketHandler[dtos.SubscribeMessageDto](app.config.WebURL)
	wsh.SetOnCloseCallback(app.repositories.WebSockets.RemoveSubscriber)

	wsh.AddSubjectHandler(string(models.AllLocations), app.allLocationsHandler)
	wsh.AddSubjectHandler(string(models.SingleLocation), app.singleLocationHandler)

	return wsh.GetHandler()
}

func (app *application) allLocationsHandler(
	w http.ResponseWriter,
	r *http.Request,
	conn *websocket.Conn,
	msg dtos.SubscribeMessageDto,
) {
	app.repositories.WebSockets.AddSubscriber(conn, msg.Subject, msg.NormalizedName)

	locationUpdateEventDtos, _ := app.getAllCurrentLocationStates(r.Context())

	err := wsjson.Write(r.Context(), conn, locationUpdateEventDtos)
	if err != nil {
		http_tools.WSErrorResponse(w, r, conn, app.repositories.WebSockets.RemoveSubscriber, err)
		return
	}

	for {
		updateEvents := app.repositories.WebSockets.GetAllUpdateEvents(conn)
		if len(updateEvents) > 0 {
			err = wsjson.Write(r.Context(), conn, updateEvents)
			if err != nil {
				http_tools.WSErrorResponse(w, r, conn, app.repositories.WebSockets.RemoveSubscriber, err)
				return
			}
		}

		if app.config.Env != config.TestEnv {
			// todo: use an events based system?
			time.Sleep(30 * time.Second) //nolint:gomnd //no magic number
		}
	}
}

func (app *application) singleLocationHandler(
	w http.ResponseWriter,
	r *http.Request,
	conn *websocket.Conn,
	msg dtos.SubscribeMessageDto,
) {
	app.repositories.WebSockets.AddSubscriber(conn, msg.Subject, msg.NormalizedName)

	for {
		updateEvent := app.repositories.WebSockets.GetByNormalizedName(conn)
		if updateEvent.NormalizedName == msg.NormalizedName {
			err := wsjson.Write(r.Context(), conn, updateEvent)
			if err != nil {
				http_tools.WSErrorResponse(w, r, conn, app.repositories.WebSockets.RemoveSubscriber, err)
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
	locations, err := app.repositories.Locations.GetAll(ctx)
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

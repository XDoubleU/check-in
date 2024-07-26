package main

import (
	"net/http"

	"github.com/xdoubleu/essentia/pkg/wstools"

	"check-in/api/internal/dtos"
)

// @Summary	WebSocket for receiving update events
// @Tags		websocket
// @Param		subscribeMessageDto	body		SubscribeMessageDto	true	"SubscribeMessageDto"
// @Success	200					{object}	LocationUpdateEvent
// @Router		/ws [get].
func (app *application) websocketsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", app.getWebSocketHandler())
}

func (app *application) getWebSocketHandler() http.HandlerFunc {
	wsh := wstools.CreateWebSocketHandler[dtos.SubscribeMessageDto](1, 100, []string{app.config.WebURL})

	// todo figure out topics + where to make them accessible (probably WS Service)
	wsh.AddTopic()

	return wsh.Handler()
}

/* todo remove
func (app *application) allLocationsHandler(
	w http.ResponseWriter,
	r *http.Request,
	conn *websocket.Conn,
	msg dtos.SubscribeMessageDto,
) {
	app.services.WebSockets.AddSubscriber(conn, msg.Subject, msg.NormalizedName)

	locationUpdateEventDtos, _ := app.getAllCurrentLocationStates(r.Context())

	err := wsjson.Write(r.Context(), conn, locationUpdateEventDtos)
	if err != nil {
		wstools.WSErrorResponse(
			w,
			r,
			conn,
			app.services.WebSockets.RemoveSubscriber,
			err,
		)
		return
	}

	for {
		updateEvents := app.services.WebSockets.GetAllUpdateEvents(conn)
		if len(updateEvents) > 0 {
			err = wsjson.Write(r.Context(), conn, updateEvents)
			if err != nil {
				httptools.WSErrorResponse(
					w,
					r,
					conn,
					app.services.WebSockets.RemoveSubscriber,
					err,
				)
				return
			}
		}

		if app.config.Env != config.TestEnv {
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
	app.services.WebSockets.AddSubscriber(conn, msg.Subject, msg.NormalizedName)

	for {
		updateEvent := app.services.WebSockets.GetByNormalizedName(conn)
		if updateEvent.NormalizedName == msg.NormalizedName {
			err := wsjson.Write(r.Context(), conn, updateEvent)
			if err != nil {
				httptools.WSErrorResponse(
					w,
					r,
					conn,
					app.services.WebSockets.RemoveSubscriber,
					err,
				)
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
*/

package main

import (
	"net/http"
)

// @Summary	WebSocket for receiving update events
// @Tags		websocket
// @Param		subscribeMessageDto	body		SubscribeMessageDto	true	"SubscribeMessageDto"
// @Success	200					{object}	LocationUpdateEvent
// @Router		/ws [get].
func (app *Application) websocketsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", app.services.WebSocket.Handler())
}

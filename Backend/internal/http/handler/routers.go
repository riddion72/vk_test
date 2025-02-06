package handler

import (
	"net/http"
)

func (h *Handler) Route() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /addresses", h.CreateAddres)
	router.HandleFunc("/addresses", h.GetAddres)

	// router.Handle("POST /house/{id}/subscribe", h.jwtMiddleware(http.HandlerFunc(h.NewSubscription), []string{"client", "moderator"}))

	return router
}

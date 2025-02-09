package handler

import (
	"net/http"
)

func (h *Handler) Route() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/put_address", h.CreateAddres)
	router.HandleFunc("/addresses", h.GetAddres)

	return router
}

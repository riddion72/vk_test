package handler

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (h *Handler) Route() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", h.Registration)
	router.HandleFunc("POST /login", h.Login)
	router.HandleFunc("GET /dummyLogin", h.DummyLogin)

	router.Handle("POST /house/{id}/subscribe", h.jwtMiddleware(http.HandlerFunc(h.NewSubscription), []string{"client", "moderator"}))
	router.Handle("POST /house/create", h.jwtMiddleware(http.HandlerFunc(h.CreateHouse), []string{"moderator"}))
	router.Handle("GET /house/{id}", h.jwtMiddleware(http.HandlerFunc(h.GetHouse), []string{"client", "moderator"}))

	router.Handle("POST /flat/create", h.jwtMiddleware(http.HandlerFunc(h.CreateFlat), []string{"client", "moderator"}))
	router.Handle("POST /flat/update", h.jwtMiddleware(http.HandlerFunc(h.UpdateFlat), []string{"moderator"}))

	return router
}

var routes = Routes{

	Route{
		"Index",
		"GET",
		"/",
		index,
	},

	Route{
		"Admin_panel",
		"GET",
		"/admin_panel",
		admin_panel,
	},

	Route{
		"Admin_panel_set",
		"POST",
		"/admin_panel",
		admin_panel_post,
	},

	Route{
		"Login",
		"GET",
		"/login",
		login,
	},
}

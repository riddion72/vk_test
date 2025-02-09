package handler

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"backend/internal/http/handler/model"
	"backend/internal/http/handler/tools"
	"backend/internal/logger"
)

const (
	PATH_TO_INDEX string = "web/teamplate/index.html"
	// PATH_TO_AUTHO string = "web/teamplate/authorization.html"
	// PATH_TO_ADMIN string = "web/teamplate/admin.html"
)

// var jwtKey = []byte("jwt_shit")

func (h *Handler) GetAddres(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetAddres"

	var pingReq model.GetAddressListRequest
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		logger.Error("page request",
			slog.String("op", op),
			slog.String("error", err.Error()))
		http.Redirect(w, r, "/addresses?page=1", http.StatusSeeOther)
		return
	}

	pingReq.Page = page

	res, err := h.pingService.GetPing(r.Context(), pingReq)
	if err != nil {
		logger.Error("get Addres",
			slog.String("op", op),
			slog.String("error", err.Error()))
		tools.SendInternalError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.New("index.html").Funcs(
		template.FuncMap{
			"sum": sum,
			"sub": sub,
		},
	)

	tmpl = template.Must(tmpl.ParseFiles(PATH_TO_INDEX))

	err = tmpl.Execute(w, *res)
	if err != nil {
		log.Println("Error execute tmpl: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateAddres(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CreateAddres"

	houseReq, err := tools.Decode[model.SetAddressListRequest](r)
	if err != nil {
		logger.Error("decode request",
			slog.String("op", op),
			slog.String("error", err.Error()))
		tools.SendClientError(w, "invalid json", http.StatusBadRequest)
		return
	}

	for _, addres := range houseReq.Addresses {
		_, err := h.pingService.CreatePing(r.Context(), addres)
		if err != nil {
			logger.Error("create ping",
				slog.String("op", op),
				slog.String("error", err.Error()))
			tools.SendInternalError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tools.SendResponse(w, http.StatusOK, http.StatusOK)
	}
}

func sum(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"main/internal/models"
	rep "main/internal/repository"
	"main/internal/usecase"
	"main/pkg/postgres"

	"github.com/golang-jwt/jwt/v5"
)

const (
	PATH_TO_INDEX string = "web/teamplate/index.html"
	PATH_TO_AUTHO string = "web/teamplate/authorization.html"
	PATH_TO_ADMIN string = "web/teamplate/admin.html"
)

var jwtKey = []byte("jwt_shit")

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	db, err := postgres.ConnectionDB()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(http.StatusInternalServerError)
		if err != nil {
			log.Println(err)
		}
	}
	defer db.Close()

	pgPepo := rep.NewPGRepository(db)
	dataBase := usecase.NewUsecase(pgPepo)

	var req models.GetArticleRequest
	pageStr := r.URL.Query().Get("page")
	if req.Page, err = strconv.Atoi(pageStr); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/?page=1", http.StatusSeeOther)
		return
	}

	res, err := dataBase.GetArticle(req)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/?page=1", http.StatusSeeOther)
		return
	}

	tmpl := template.New("index.html").Funcs(
		template.FuncMap{
			"sum": sum,
			"sub": sub,
		},
	)

	tmpl = template.Must(tmpl.ParseFiles(PATH_TO_INDEX))
	// log.Println(res)

	err = tmpl.Execute(w, res)
	if err != nil {
		log.Println("Error execute tmpl: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if usecase.Authorization(r.FormValue("login"), r.FormValue("password")) {
		login := r.FormValue("login")
		expirationDate := time.Now().Add(1 * time.Hour)
		claims := jwt.MapClaims{
			"exp":   expirationDate.Unix(),
			"login": login,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(jwtKey)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "JWT erorr", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenStr,
			Expires: expirationDate,
		})

		http.Redirect(w, r, "/admin_panel", http.StatusSeeOther)
	}

	tmpl := template.Must(template.ParseFiles(PATH_TO_AUTHO))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error execute tmpl: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	return
}

func (h *Handler) admin_panel(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("No token")
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}
		log.Println(err.Error())
		http.Error(w, "No Cookie?", http.StatusBadRequest)
		return
	}

	token := c.Value
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("Wrong token")
			http.Error(w, "Wrong token1", http.StatusUnauthorized)
			return
		}
		log.Println(err.Error())
		http.Error(w, "Wrong token2", http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		http.Error(w, "Wrong token3", http.StatusUnauthorized)
		return
	}

	http.ServeFile(w, r, PATH_TO_ADMIN)
}

func (h *Handler) admin_panel_post(w http.ResponseWriter, r *http.Request) {

	// if r.Method == "POST" {
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	db, err := postgres.ConnectionDB()
	if err != nil {
		log.Println("Erorr ConnectionDB: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Erorr ConnectionDB", http.StatusInternalServerError)
	}
	defer db.Close()

	pgPepo := rep.NewPGRepository(db)
	dataBase := usecase.NewUsecase(pgPepo)

	var article models.SetArticleRequest
	article.Request.Article_name = r.Form.Get("title")
	article.Request.Article_content = r.Form.Get("content")
	err = dataBase.SetArticle(article)
	if err != nil {
		log.Println("Erorr SetArticle: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Erorr SetArticle", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, PATH_TO_ADMIN)
	// }
}

func sum(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

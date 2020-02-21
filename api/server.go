package api

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/halaalajlan/hackathon/auth"
	ctx "github.com/halaalajlan/hackathon/context"
	log "github.com/halaalajlan/hackathon/logger"
	mid "github.com/halaalajlan/hackathon/middleware"
	"github.com/halaalajlan/hackathon/models"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	defaultServer := &http.Server{
		ReadTimeout: 10 * time.Second,
		Addr:        "127.0.0.1:3333",
	}
	as := &Server{server: defaultServer}
	as.registerRoutes()
	return as
}
func (as *Server) Start() {

	log.Infof("Starting  server at http://%s", "127.0.0.1:3333")
	log.Fatal(as.server.ListenAndServe())
}

// Shutdown attempts to gracefully shutdown the server.
func (as *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return as.server.Shutdown(ctx)
}

func (as *Server) registerRoutes() {
	router := mux.NewRouter()
	router.StrictSlash(false)
	router.Use(mid.GetContext)
	router.Use(mid.RequireAPIKey)
	router.HandleFunc("/", mid.Use(as.Base, mid.RequireLogin)) //redirect user to home page
	router.HandleFunc("/api/login", as.Login)
	router.HandleFunc("/login", as.LoginUI)
	router.HandleFunc("/home", mid.Use(as.Base, mid.RequireLogin))

	router.HandleFunc("/api/patient", as.Patients)
	handler := handlers.CombinedLoggingHandler(log.Writer(), router)
	as.server.Handler = handler
}

func (as *Server) LoginUI(w http.ResponseWriter, r *http.Request) {
	params := struct {
		User    models.Hospital
		Title   string
		Flashes []interface{}
	}{Title: "Login"}
	session := ctx.Get(r, "session").(*sessions.Session)
	switch {
	case r.Method == "GET":
		params.Flashes = session.Flashes()
		session.Save(r, w)
		templates := template.New("template")
		_, err := templates.ParseFiles("templates/login.html", "templates/flashes.html", "templates/header.html")
		if err != nil {
			log.Error(err)
		}
		template.Must(templates, err).ExecuteTemplate(w, "base", params)
		return
	case r.Method == "POST":
		//Attempt to login
		succ, u, err := auth.Login(r)
		if err != nil {
			log.Error(err)
		}
		//If we've logged in, save the session and redirect to the home page
		if succ {
			session.Values["id"] = u.Id
			session.Save(r, w)
			next := "/"
			url, err := url.Parse(r.FormValue("next"))
			if err == nil {
				path := url.Path
				if path != "" {
					next = path
				}
			}
			http.Redirect(w, r, next, 302)
		} else {
			Flash(w, r, "danger", "Invalid Username/Password")
			params.Flashes = session.Flashes()
			session.Save(r, w)
			templates := template.New("template")
			_, err := templates.ParseFiles("templates/login.html", "templates/flashes.html", "templates/header.html")
			if err != nil {
				log.Error(err)
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			template.Must(templates, err).ExecuteTemplate(w, "base", params)
		}
	}
}

// Base handles the default path and template execution
func (as *Server) Base(w http.ResponseWriter, r *http.Request) {
	templates := template.New("template")
	_, err := templates.ParseFiles("templates/index.html", "templates/flashes.html", "templates/header.html")
	if err != nil {
		log.Error(err)
	}
	template.Must(templates, err).ExecuteTemplate(w, "base", nil)
}

// Flash handles the rendering flash messages
func Flash(w http.ResponseWriter, r *http.Request, t string, m string) {
	session := ctx.Get(r, "session").(*sessions.Session)
	session.AddFlash(models.Flash{
		Type:    t,
		Message: m,
	})
}

package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, method string, handler http.HandlerFunc) {
	if s.Handlers[path] == nil {
		s.Handlers[path] = make(map[string]http.HandlerFunc)
	}
	s.Handlers[path][method] = handler
}

// Loop through the handlers and add them to the router
// register middleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, methods := range s.Handlers {
		for method, handler := range methods {
			switch method {
			case "GET":
				s.Router.Get(path, handler)
			case "POST":
				s.Router.Post(path, handler)
			case "PUT":
				s.Router.Put(path, handler)
			case "DELETE":
				s.Router.Delete(path, handler)
			}
		}
	}
	http.ListenAndServe(":"+s.WebServerPort, s.Router)
}

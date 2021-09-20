package apiserver

import (
	"encoding/json"
	"net/http"
	"pingram/pkg/apiserver/views"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
}

func newServer() *server {
	s := &server{
		router: mux.NewRouter(),
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()
	v1Router.HandleFunc("/all-views", s.handleGetAddDomains()).Methods("GET")
	v1Router.HandleFunc("/http-errors", s.handleGetHttpErrorsDomains()).Methods("GET")
	v1Router.HandleFunc("/currently-down", s.handleGetCurrentlyDownedDomains()).Methods("GET")
	v1Router.HandleFunc("/currently-month-downtimes-at-work", s.handleGetCurrentlyMonthDownedDomains()).Methods("GET")
	//v1Router.HandleFunc("/currently-month-downtimes-total", s.handleGetCurrentlyMonthDownedDomains()).Methods("GET")
	v1Router.HandleFunc("/previous-month-downtimes-at-work", s.handleGetPreviousMonthDownedDomains()).Methods("GET")
	//v1Router.HandleFunc("/previous-month-downtimes-total", s.handleGetPreviousMonthDownedDomains()).Methods("GET")
}

func (s *server) handleGetAddDomains() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allDomains := views.GetAllDomains()
		s.respond(w, http.StatusOK, allDomains)
	}
}

func (s *server) handleGetHttpErrorsDomains() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		domains := views.GetInfoHttpErrors()
		s.respond(w, http.StatusOK, domains)
	}
}

func (s *server) handleGetCurrentlyDownedDomains() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		failedDomains := views.GetDowned()
		s.respond(w, http.StatusOK, failedDomains)
	}
}

func (s *server) handleGetCurrentlyMonthDownedDomains() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		failedDomains := views.GetDownedAtCurrentMonth()
		s.respond(w, http.StatusOK, failedDomains)
	}
}

func (s *server) handleGetPreviousMonthDownedDomains() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		failedDomains := views.GetDownedAtPreviousMonth()
		s.respond(w, http.StatusOK, failedDomains)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (s *server) error(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	s.respond(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

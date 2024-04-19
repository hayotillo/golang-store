package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"store-api/app/service"
	"store-api/app/store"
)

type server struct {
	router       *mux.Router
	store        store.Store
	service      service.Service
	jwtSecretKey string
}

func NewServer(store store.Store, service service.Service, jwtSecretKey string) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		service:      service,
		jwtSecretKey: jwtSecretKey,
	}
	s.init()
	return s
}

func (s *server) init() {
	// set global middlewares
	router := s.router.NewRoute().Subrouter()
	router.Use(s.JsonMiddleware)
	// public
	router.HandleFunc("/user/login", s.handleUserLogin()).Methods("POST")
	// private
	private := router.PathPrefix("/private").Subrouter()
	private.Use(s.UserAuthenticatedMiddleware)
	// user
	user := private.PathPrefix("/user").Subrouter()
	user.HandleFunc("/get-by-token", s.handleUserByToken()).Methods("POST")
	user.HandleFunc("/logout", s.handleUserLogout()).Methods("POST")
	user.HandleFunc("/one", s.handleUserOne()).Methods("POST")
	user.HandleFunc("/list", s.handleUserList()).Methods("POST")
	user.HandleFunc("/save", s.handleUserSave()).Methods("POST")
	user.HandleFunc("/delete", s.handleUserDelete()).Methods("POST")
	// visit
	visit := private.PathPrefix("/visit").Subrouter()
	visit.HandleFunc("/one", s.handleVisitOne()).Methods("POST")
	visit.HandleFunc("/list", s.handleVisitList()).Methods("POST")
	visit.HandleFunc("/save", s.handleVisitSave()).Methods("POST")
	visit.HandleFunc("/delete", s.handleVisitDelete()).Methods("POST")
	// customer
	customer := private.PathPrefix("/customer").Subrouter()
	customer.HandleFunc("/get", s.handleCustomerGet()).Methods("POST")
	customer.HandleFunc("/one", s.handleCustomerOne()).Methods("POST")
	customer.HandleFunc("/list", s.handleCustomerList()).Methods("POST")
	customer.HandleFunc("/save", s.handleCustomerSave()).Methods("POST")
	customer.HandleFunc("/delete", s.handleCustomerDelete()).Methods("POST")
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
	}
	s.respond(w, r, code, map[string]string{"error": errMessage})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

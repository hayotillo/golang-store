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
	//router.HandleFunc("/check-file", s.handleSaleCheckFile()).Methods("GET")
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
	// sale
	sale := private.PathPrefix("/sale").Subrouter()
	sale.HandleFunc("/history", s.handleSaleHistory()).Methods("POST")
	sale.HandleFunc("/check-file", s.handleSaleCheckFile()).Methods("POST")
	sale.HandleFunc("/get", s.handleSaleGet()).Methods("POST")
	sale.HandleFunc("/one", s.handleSaleOne()).Methods("POST")
	sale.HandleFunc("/list", s.handleSaleList()).Methods("POST")
	sale.HandleFunc("/products", s.handleSaleProducts()).Methods("POST")
	sale.HandleFunc("/save", s.handleSaleSave()).Methods("POST")
	sale.HandleFunc("/delete", s.handleSaleDelete()).Methods("POST")
	// product
	product := private.PathPrefix("/product").Subrouter()
	product.HandleFunc("/one", s.handleProductOne()).Methods("POST")
	product.HandleFunc("/list", s.handleProductList()).Methods("POST")
	product.HandleFunc("/save", s.handleProductSave()).Methods("POST")
	product.HandleFunc("/delete", s.handleProductDelete()).Methods("POST")
	// incoming
	incoming := private.PathPrefix("/incoming").Subrouter()
	incoming.HandleFunc("/get", s.handleIncomingGet()).Methods("POST")
	incoming.HandleFunc("/one", s.handleIncomingOne()).Methods("POST")
	incoming.HandleFunc("/list", s.handleIncomingList()).Methods("POST")
	incoming.HandleFunc("/save", s.handleIncomingSave()).Methods("POST")
	incoming.HandleFunc("/delete", s.handleIncomingDelete()).Methods("POST")

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

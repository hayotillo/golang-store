package server

import (
	"github.com/gorilla/schema"
	"net/http"
	"reception/app/model"
)

func (s *server) handleCustomerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.CustomerOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		m, err := s.store.Customer().Get(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleCustomerOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.CustomerOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		m, err := s.store.Customer().One(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleCustomerList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.CustomerListFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		list, err := s.store.Customer().List(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, list)
	}
}

func (s *server) handleCustomerSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		m := &model.CustomerData{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(m, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.Customer().Save(m)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleCustomerDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.CustomerDeleteFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.Customer().Delete(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]bool{"result": true})
	}
}

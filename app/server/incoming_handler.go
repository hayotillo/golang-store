package server

import (
	"github.com/gorilla/schema"
	"net/http"
	"store-api/app/model"
)

func (s *server) handleIncomingGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.IncomingOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		m, err := s.store.Incoming().Get(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleIncomingOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.IncomingOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		m, err := s.store.Incoming().One(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleIncomingList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.IncomingListFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		list, err := s.store.Incoming().List(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, list)
	}
}

func (s *server) handleIncomingSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		m := &model.IncomingData{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(m, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := r.Context().Value(s.jwtSecretKey).(*model.User)
		m.UserID = u.ID

		err = s.store.Incoming().Save(m)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleIncomingDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.IncomingDeleteFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.Incoming().Delete(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]bool{"result": true})
	}
}

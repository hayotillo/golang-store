package server

import (
	"github.com/gorilla/schema"
	"net/http"
	"store-api/app/model"
)

func (s *server) handleProductOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.ProductOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		m, err := s.store.Product().One(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleProductList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.ProductListFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		list, err := s.store.Product().List(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, list)
	}
}

func (s *server) handleProductSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		m := &model.Product{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(m, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.Product().Save(m)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleProductDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.ProductDeleteFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.Product().Delete(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]bool{"result": true})
	}
}

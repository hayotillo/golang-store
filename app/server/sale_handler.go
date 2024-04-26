package server

import (
	"fmt"
	"github.com/gorilla/schema"
	"net/http"
	"store-api/app/model"
)

func (s *server) handleSaleHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.SaleHistoryFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		m, err := s.store.Sale().History(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleSaleCheckFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.SaleOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		checkFile, err := s.store.Sale().CheckFile(f)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s1.pdf", checkFile.Name))

		if checkFile == nil {
			s.error(w, r, http.StatusUnprocessableEntity, nil)
			return
		}

		w.Write(checkFile.File)
		return
	}
}

func (s *server) handleSaleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.SaleOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		m, err := s.store.Sale().Get(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleSaleOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.SaleOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		m, err := s.store.Sale().One(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleSaleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.SaleListFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		list, err := s.store.Sale().List(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, list)
	}
}

func (s *server) handleSaleProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.SaleProductListFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		list, err := s.store.Sale().Products(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, list)
	}
}

func (s *server) handleSaleSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		m := &model.SaleData{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(m, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := r.Context().Value(s.jwtSecretKey).(*model.User)
		m.UserID = u.ID

		err = s.store.Sale().Save(m)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleSaleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		f := model.SaleDeleteFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.Sale().Delete(f)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]bool{"result": true})
	}
}

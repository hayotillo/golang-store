package server

import (
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/schema"
	"net/http"
	"reception/app/model"
	"reception/app/store"
	"strings"
)

func (s *server) handleUserGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		f := model.UserOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().One(f)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		s.respond(w, r, http.StatusOK, u.ToPublic())
	}
}

func (s *server) handleUserOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(s.jwtSecretKey).(*model.User)
		//if !u.HasCompanyPermission() {
		//	s.error(w, r, http.StatusBadRequest, store.ErrUserPermissionDenied)
		//	return
		//}
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		f := model.UserOneFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if !u.IsAdmin() && u.ID != f.ID {
			f.UserID = u.ID
		}

		m, err := s.store.User().One(f)
		//fmt.Println("one", err)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleUserRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		data := model.User{}
		decoder := schema.NewDecoder()

		if err := decoder.Decode(&data, r.PostForm); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if len(data.Phone) < 4 || len(data.Password) < 4 {
			s.error(w, r, http.StatusBadRequest, store.ErrRequiredDataNotFount)
			return
		}
		exists, err := s.store.User().GetByPhone(data.PhoneData)
		if exists != nil && exists.CheckIDData() {
			s.error(w, r, http.StatusAlreadyReported, store.ErrRecordAlreadyExists)
			return
		}
		data.Status = "user"
		m, err := s.store.User().Save(data)

		if err != nil {
			s.error(w, r, http.StatusAlreadyReported, err)
			return
		}
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.MapClaims{"id": m.ID},
		)
		secretToken, err := token.SignedString([]byte(s.jwtSecretKey))
		ok := s.store.User().SetToken(m.IDData, secretToken)
		if ok {
			m.Token = secretToken
		}
		s.respond(w, r, http.StatusCreated, m)
	}
}

func (s *server) handleUserSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(1 << 20)
		if err != nil {
			err = r.ParseForm()
		}
		data := model.User{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&data, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(s.jwtSecretKey).(*model.User)
		if !data.CheckIDData() && !data.CheckUserStatusData() {
			s.error(w, r, http.StatusBadRequest, store.ErrRequiredDataNotFount)
			return
		}
		if !u.IsAdmin() && data.Status == "admin" {
			data.Status = ""
		}

		m, err := s.store.User().Save(data)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleUserList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(s.jwtSecretKey).(*model.User)
		if !u.IsAdmin() {
			s.error(w, r, http.StatusBadRequest, store.ErrUserPermissionDenied)
			return
		}

		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		f := model.UserListFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		list, err := s.store.User().List(f)
		if err != nil {
			s.error(w, r, http.StatusBadGateway, err)
			return
		}

		s.respond(w, r, http.StatusOK, list)
	}
}

func (s *server) handleUserByToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(s.jwtSecretKey).(*model.User)
		one, err := s.store.User().One(model.UserOneFilter{IDData: model.IDData{ID: u.ID}})
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}
		s.respond(w, r, http.StatusOK, one)
	}
}

func (s *server) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		f := model.UserDeleteFilter{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&f, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := r.Context().Value(s.jwtSecretKey).(*model.User)

		if !u.IsAdmin() && u.ID != f.UserID {
			s.error(w, r, http.StatusBadRequest, store.ErrUserPermissionDenied)
			return
		}
		if !u.IsAdmin() || !f.CheckUserIDData() {
			f.UserID = u.ID
		}

		res, err := s.store.User().Delete(f)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]bool{"result": res})
	}
}

func (s *server) handleUserSetStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(s.jwtSecretKey).(*model.User)
		if !u.IsAdmin() {
			s.error(w, r, http.StatusForbidden, store.ErrUserPermissionDenied)
			return
		}
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		f := model.UserStatusFilterData{}
		decoder := schema.NewDecoder()
		if err = decoder.Decode(&f, r.PostForm); err == nil {
			res, err := s.store.User().SetStatus(f)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			s.respond(w, r, http.StatusOK, map[string]bool{"result": res})
		}
	}
}

func (s *server) handleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		data := model.UserLoginData{}
		decoder := schema.NewDecoder()
		err = decoder.Decode(&data, r.PostForm)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		m, err := s.store.User().GetByPhone(data.PhoneData)
		if err != nil || !m.ComparePassword(data.Password) {
			s.error(w, r, http.StatusUnauthorized, store.ErrLoginOrPasswordIncorrect)
			return
		}

		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.MapClaims{"id": m.ID},
		)
		secretToken, err := token.SignedString([]byte(s.jwtSecretKey))

		if m != nil && err == nil {
			s.store.User().SetToken(m.IDData, secretToken)
			m.Token = secretToken
			s.respond(w, r, http.StatusOK, m)
		} else {
			s.respond(w, r, http.StatusUnauthorized, nil)
		}
	}
}

func (s *server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			token := strings.TrimPrefix(authorization, "Token ")
			if len(token) > 0 {
				claims := jwt.MapClaims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(s.jwtSecretKey), nil
				})
				if err == nil {
					for k, v := range claims {
						if k == "id" {
							u, err := s.store.User().One(model.UserOneFilter{IDData: model.IDData{
								ID: v.(string),
							}})
							if err == nil && u != nil {
								if u.Token == token {
									rem := s.store.User().SetToken(u.IDData, "")
									s.respond(w, r, http.StatusOK, map[string]bool{"result": rem})
									return
								}
							}
							break
						}
					}
					return
				}
			}
		}
		s.error(w, r, http.StatusUnauthorized, store.ErrUserNotAuthenticated)
		return
	}
}

package store

import "reception/app/model"

type UserInterface interface {
	GetByPhone(p model.PhoneData) (*model.User, error)
	One(f model.UserOneFilter) (*model.User, error)
	List(filter model.UserListFilter) (*model.ListData, error)
	Save(m model.User) (*model.User, error)
	Delete(f model.UserDeleteFilter) (bool, error)
	SetToken(d model.IDData, t string) bool
	SetStatus(d model.UserStatusFilterData) (bool, error)
}

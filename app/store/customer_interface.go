package store

import "reception/app/model"

type CustomerInterface interface {
	Get(f model.CustomerOneFilter) (*model.Customer, error)
	One(f model.CustomerOneFilter) (*model.CustomerData, error)
	List(f model.CustomerListFilter) (*model.ListData, error)
	Save(m *model.CustomerData) error
	Delete(f model.CustomerDeleteFilter) error
}

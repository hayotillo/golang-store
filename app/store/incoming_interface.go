package store

import "store-api/app/model"

type IncomingInterface interface {
	Get(f model.IncomingOneFilter) (*model.Incoming, error)
	One(f model.IncomingOneFilter) (*model.IncomingData, error)
	List(f model.IncomingListFilter) (*model.ListData, error)
	Save(m *model.IncomingData) error
	Delete(f model.IncomingDeleteFilter) error
}

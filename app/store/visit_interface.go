package store

import "store-api/app/model"

type VisitInterface interface {
	Get(f model.VisitOneFilter) (*model.Visit, error)
	One(f model.VisitOneFilter) (*model.VisitData, error)
	List(f model.VisitListFilter) (*model.ListData, error)
	Save(m *model.VisitData) error
	Delete(f model.VisitDeleteFilter) error
}

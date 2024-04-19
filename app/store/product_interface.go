package store

import "store-api/app/model"

type ProductInterface interface {
	One(f model.ProductOneFilter) (*model.Product, error)
	List(f model.ProductListFilter) (*model.ListData, error)
	Save(m *model.Product) error
	Delete(f model.ProductDeleteFilter) error
}

package store

import "store-api/app/model"

type SaleInterface interface {
	CheckFile(f model.SaleOneFilter) (*model.ReportFileData, error)
	Get(f model.SaleOneFilter) (*model.Sale, error)
	One(f model.SaleOneFilter) (*model.SaleData, error)
	List(f model.SaleListFilter) (*model.ListData, error)
	Products(f model.SaleProductListFilter) (*model.ListData, error)
	Save(m *model.SaleData) error
	Delete(f model.SaleDeleteFilter) error
}

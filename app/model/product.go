package model

type Product struct {
	IDData
	NameData
	IncomingSum string `json:"incoming_sum" schema:"incoming_sum"`
	SaleSum     string `json:"sale_sum" schema:"sale_sum"`
}

func (d *Product) CheckInsertData() bool {
	return d.CheckNameData()
}

func (d *Product) CheckUpdateData() bool {
	return d.CheckIDData() && d.CheckNameData()
}

type ProductOneFilter struct {
	IDData
}

type ProductListFilter struct {
	SearchData
	OrdersData
	PaginateData
}

type ProductDeleteFilter struct {
	IDSData
}

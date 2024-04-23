package model

type Product struct {
	IDData
	NameData
	IncomingCount string `json:"incoming_count" schema:"incoming_count"`
	SaleCount     string `json:"sale_count" schema:"sale_count"`
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

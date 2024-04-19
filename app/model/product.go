package model

type Product struct {
	IDData
	NameData
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
	PaginateData
}

type ProductDeleteFilter struct {
	IDSData
}

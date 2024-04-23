package model

type Incoming struct {
	IDData
	PriceData
	QuantityData
	ArchiveData
	DescriptionData
	TimestampData
	User    User    `json:"user" schema:"-"`
	Product Product `json:"product" schema:"-"`
}

type IncomingData struct {
	IDData
	UserIDData
	ProductIDData
	PriceData
	QuantityData
	ArchiveData
	DescriptionData
	TimestampData
}

func (d *IncomingData) CheckInsertData() bool {
	return d.CheckUserIDData() &&
		d.CheckProductIDData() &&
		d.CheckPriceData() &&
		d.CheckQuantityData()
}

func (d *IncomingData) CheckUpdateData() bool {
	return d.CheckIDData() && (d.CheckArchiveData() ||
		d.CheckUserIDData() ||
		d.CheckProductIDData() ||
		d.CheckPriceData() ||
		d.CheckQuantityData() ||
		d.CheckDescriptionData())
}

type IncomingOneFilter struct {
	IDData
}

type IncomingListFilter struct {
	SearchData
	OrdersData
	ProductIDData
	ArchiveData
	PaginateData
}

type IncomingDeleteFilter struct {
	IDSData
}

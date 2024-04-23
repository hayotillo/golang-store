package model

type Incoming struct {
	IDData
	PriceData
	QuantityData
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
	return d.CheckIDData() && (d.CheckUserIDData() ||
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
	PaginateData
}

type IncomingDeleteFilter struct {
	IDSData
}

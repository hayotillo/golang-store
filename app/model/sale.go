package model

type Sale struct {
	IDData
	DescriptionData
	ArchiveData
	TimestampData
	PriceData
	User    User      `json:"user" schema:"-"`
	Product []Product `json:"products,omitempty" schema:"-"`
}

type SaleData struct {
	IDData
	UserIDData
	SaleOrdersData
	DescriptionData
	ArchiveData
	ClearOrdersData
	TimestampData
}

func (d *SaleData) CheckInsertData() bool {
	return d.CheckUserIDData() && d.CheckSaleOrdersData()
}

func (d *SaleData) CheckUpdateData() bool {
	return d.CheckIDData() && (d.CheckUserIDData() ||
		d.CheckClearOrdersData() ||
		d.CheckDescriptionData() ||
		d.CheckArchiveData())
}

type ClearOrdersData struct {
	ClearOrders string `json:"-" schema:"clear_orders"`
}

func (d *ClearOrdersData) CheckClearOrdersData() bool {
	return d.ClearOrders == "true" || d.ClearOrders == "false"
}

func (d *ClearOrdersData) IsClearOrders() bool {
	return d.ClearOrders == "true"
}

type DescriptionData struct {
	Description string `json:"description" schema:"description"`
}

func (d *DescriptionData) CheckDescriptionData() bool {
	return len(d.Description) > 0
}

type ArchiveData struct {
	Archive string `json:"archive" schema:"archive"`
}

func (d *ArchiveData) CheckArchiveData() bool {
	return d.Archive == "true" || d.Archive == "false"
}

func (d *ArchiveData) IsArchive() bool {
	return d.Archive == "true"
}

type SaleProduct struct {
	IDData
	NameData
	PriceData
	QuantityData
}

type SaleOneFilter struct {
	IDData
}

type SaleListFilter struct {
	SearchData
	PaginateData
}

type SaleProductListFilter struct {
	SearchData
	SaleIDData
	OrdersData
	PaginateData
}

type SaleDeleteFilter struct {
	IDSData
}

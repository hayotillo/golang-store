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

type SelfOnlyData struct {
	SelfOnly string `json:"-" schema:"self_only"`
}

func (d *SelfOnlyData) CheckSelfOnlyData() bool {
	return d.SelfOnly == "true" || d.SelfOnly == "false"
}

func (d *SelfOnlyData) IsSelfOnly() bool {
	return d.SelfOnly == "true"
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

type SaleHistories struct {
	Items []SaleHistoryItem `json:"items" schema:"items"`
}

type SaleHistoryItem struct {
	IDData
	PriceData
	QuantityData
	TimestampData
	User User `json:"user" schema:"user"`
}

type SaleOneFilter struct {
	IDData
}

type SaleListFilter struct {
	SearchData
	PeriodData
	PaginateData
}

type SaleProductListFilter struct {
	SearchData
	SaleIDData
	OrdersData
	SelfOnlyData
	PaginateData
}

type SaleHistoryFilter struct {
	SaleIDData
	OrdersData
	SelfOnlyData
	PeriodData
}

type SaleDeleteFilter struct {
	IDSData
}

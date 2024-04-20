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
	TimestampData
}

func (d *SaleData) CheckInsertData() bool {
	return d.CheckUserIDData() && d.CheckSaleOrdersData()
}

func (d *SaleData) CheckUpdateData() bool {
	return d.CheckIDData() && (d.CheckUserIDData() &&
		d.CheckSaleOrdersData() &&
		d.CheckDescriptionData() &&
		d.CheckArchiveData())
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

type SaleOneFilter struct {
	IDData
}

type SaleListFilter struct {
	SearchData
	PaginateData
}

type SaleDeleteFilter struct {
	IDSData
}

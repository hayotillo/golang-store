package model

type Sale struct {
	IDData
	StatusData
	PriceData
	QuantityData
	DescriptionData
	ArchiveData
	TimestampData
	User    User    `json:"user" schema:"-"`
	Product Product `json:"product" schema:"-"`
}

type SaleData struct {
	IDData
	UserIDData
	ProductIDData
	StatusData
	PriceData
	QuantityData
	DescriptionData
	ArchiveData
	TimestampData
}

func (d *SaleData) CheckInsertData() bool {
	return d.CheckUserIDData() && d.CheckProductIDData()
}

func (d *SaleData) CheckUpdateData() bool {
	return d.CheckIDData() && (d.CheckUserIDData() &&
		d.CheckProductIDData() &&
		d.CheckUserStatusData() &&
		d.CheckDescriptionData() &&
		d.CheckArchiveData())
}

func (d *StatusData) CheckStatusData() bool {
	switch d.Status {
	case "archive", "wait":
		return true
		break
	default:
		return false
	}
	return false
}

type DescriptionData struct {
	Description string `json:"diagnosis" schema:"diagnosis"`
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

package model

type Customer struct {
	IDData
	FullNameData
	BirthData
	PhoneData
	TimestampData
	Visit VisitData `json:"visit" schema:"visit"`
}

type CustomerData struct {
	IDData
	DoctorIDData
	FullNameData
	BirthData
	PhoneData
}

func (d *CustomerData) CheckInsertData() bool {
	return d.CheckFullNameData() && d.CheckBirthData()
}

func (d *CustomerData) CheckUpdateData() bool {
	return d.CheckIDData() && (d.CheckFullNameData() ||
		d.CheckBirthData() ||
		d.CheckPhoneData())
}

type BirthData struct {
	Birth string `json:"birth,omitempty" schema:"birth"`
}

func (d *BirthData) CheckBirthData() bool {
	return len(d.Birth) == 10
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

type CustomerOneFilter struct {
	IDData
	PhoneData
}

type CustomerListFilter struct {
	SearchData
	PaginateData
}

type CustomerDeleteFilter struct {
	IDSData
}

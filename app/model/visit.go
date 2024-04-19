package model

type Visit struct {
	IDData
	StatusData
	DiagnosisData
	ArchiveData
	TimestampData
	User     User     `json:"user" schema:"user"`
	Customer Customer `json:"customer" schema:"customer"`
}

type VisitData struct {
	IDData
	UserIDData
	CustomerIDData
	StatusData
	DiagnosisData
	ArchiveData
	TimestampData
}

func (d *VisitData) CheckInsertData() bool {
	return d.CheckUserIDData() && d.CheckCustomerIDData()
}

func (d *VisitData) CheckUpdateData() bool {
	return d.CheckIDData() && (d.CheckUserIDData() &&
		d.CheckCustomerIDData() &&
		d.CheckUserStatusData() &&
		d.CheckDiagnosisData() &&
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

type DiagnosisData struct {
	Diagnosis string `json:"diagnosis" schema:"diagnosis"`
}

func (d *DiagnosisData) CheckDiagnosisData() bool {
	return len(d.Diagnosis) > 0
}

type VisitOneFilter struct {
	IDData
}

type VisitListFilter struct {
	SearchData
	PaginateData
}

type VisitDeleteFilter struct {
	IDSData
}

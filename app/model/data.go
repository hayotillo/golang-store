package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type IDData struct {
	ID string `json:"id,omitempty" schema:"id"`
}

func (d *IDData) CheckIDData() bool {
	return len(d.ID) == 36
}

type CustomerIDData struct {
	CustomerID string `json:"customer_id,omitempty" schema:"customer_id"`
}

func (d *CustomerIDData) CheckCustomerIDData() bool {
	return len(d.CustomerID) == 36
}

type DoctorIDData struct {
	DoctorID string `json:"doctor_id,omitempty" schema:"doctor_id"`
}

func (d *DoctorIDData) CheckDoctorIDData() bool {
	return len(d.DoctorID) == 36
}

type ListData struct {
	Items    []interface{} `json:"items"`
	Paginate Paginate      `json:"paginate"`
}

type IDSData struct {
	JsonIDS string `json:"-" schema:"ids"`
	ids     []string
}

func (d *IDSData) IDS() []string {
	if len(d.ids) == 0 && len(d.JsonIDS) > 3 {
		err := json.Unmarshal([]byte(d.JsonIDS), &d.ids)
		if err != nil {
			fmt.Println("error parse ids", d.JsonIDS, err)
		}
	}
	return d.ids
}

func (d *IDSData) HasIDS() bool {
	return len(d.IDS()) > 0
}

type FullNameData struct {
	FullName string `json:"full_name,omitempty" schema:"full_name"`
}

func (d *FullNameData) CheckFullNameData() bool {
	return len(d.FullName) > 1
}

type StatusData struct {
	Status string `json:"status"`
}

type TimestampData struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type OrdersData struct {
	JsonOrder string `schema:"orders"`
}

func (d *OrdersData) OrdersWhere(orders map[string]string) string {
	res := ""
	if d.CheckOrderData() {
		var data map[string]string
		err := json.Unmarshal([]byte(d.JsonOrder), &data)
		if err != nil {
			fmt.Println("error parse orders", d.JsonOrder, err)
		}

		for k, v := range data {
			orders[k] = v
		}
	}

	for f, o := range orders {
		order := strings.ToUpper(o)
		if order == "ASC" || order == "DESC" {
			res = fmt.Sprintf("%s %s %s,", res, f, order)
		}
	}

	return strings.TrimSuffix(res, ",")
}

func (d *OrdersData) CheckOrderData() bool {
	return len(d.JsonOrder) > 3
}

type UserIDData struct {
	UserID string `json:"user_id" schema:"user_id"`
}

func (d *UserIDData) CheckUserIDData() bool {
	return len(d.UserID) == 36
}

type VisitIDData struct {
	VisitID string `json:"visit_id" schema:"visit_id"`
}

func (d *VisitIDData) CheckVisitIDData() bool {
	return len(d.VisitID) == 36
}

type PhoneData struct {
	Phone string `json:"phone,omitempty" schema:"phone"`
}

func (d *PhoneData) CheckPhoneData() bool {
	l := len(d.Phone)
	return l == 13 || l == 9
}

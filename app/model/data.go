package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type IDData struct {
	ID string `json:"id,omitempty" schema:"id"`
}

func (d *IDData) CheckIDData() bool {
	return len(d.ID) == 36
}

type ProductIDData struct {
	ProductID string `json:"product_id,omitempty" schema:"product_id"`
}

func (d *ProductIDData) CheckProductIDData() bool {
	return len(d.ProductID) == 36
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
	Status string `json:"status,omitempty"`
}

type TimestampData struct {
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type OrdersData struct {
	JsonOrder string `schema:"SaleOrders"`
}

func (d *OrdersData) OrdersWhere(fields []string, def string) string {
	res := ""
	var orders []string
	if d.CheckOrderData() {
		err := json.Unmarshal([]byte(d.JsonOrder), &orders)
		if err != nil {
			fmt.Println("error parse order json", d.JsonOrder, err)
		}
	}

	for _, o := range orders {
		res = fmt.Sprintf("%s%s", res, d.parseOrder(o, fields))
	}
	if len(orders) == 0 && len(def) > 4 {
		for _, o := range strings.Split(def, ",") {
			res = fmt.Sprintf("%s%s", res, d.parseOrder(o, fields))
		}
	}

	return strings.TrimSuffix(res, ",")
}

func (d *OrdersData) parseOrder(j string, fields []string) string {
	res := ""
	s := ""
	so := strings.Split(strings.TrimSpace(j), ":")
	if len(so) == 2 {
		if strings.ToLower(strings.TrimSpace(so[1])) == "desc" {
			s = "DESC"
		} else {
			s = "ASC"
		}

		i, err := strconv.Atoi(strings.TrimSpace(so[0]))
		if err == nil {
			field := fields[i-1]
			res = fmt.Sprintf("%s %s %s,", res, field, s)
		}
	}

	return res
}

func (d *OrdersData) CheckOrderData() bool {
	return len(d.JsonOrder) > 3
}

type SaleOrdersData struct {
	SaleOrdersJson string `json:"-" schema:"sale_orders"`
	SaleOrders     []struct {
		ProductIDData
		PriceData
		QuantityData
	} `json:"-" schema:"-"`
}

func (d *SaleOrdersData) CheckSaleOrdersData() bool {
	if len(d.SaleOrders) == 0 && len(d.SaleOrdersJson) > 40 {
		err := json.Unmarshal([]byte(d.SaleOrdersJson), &d.SaleOrders)
		if err != nil {
			fmt.Println("error parse sale SaleOrders", d.SaleOrdersJson, err)
		}
	}
	return len(d.SaleOrders) > 0
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

type NameData struct {
	Name string `json:"name,omitempty" schema:"name"`
}

func (d *NameData) CheckNameData() bool {
	return len(d.Name) > 0
}

type PriceData struct {
	Price string `json:"price,omitempty" schema:"price"`
}

func (d *PriceData) CheckPriceData() bool {
	if len(d.Price) > 0 {
		_, err := strconv.ParseFloat(d.Price, 32)
		if err != nil {
			fmt.Println("price data error", err)
		}

		if err == nil {
			return true
		}
	}
	return false
}

func (d *PriceData) PriceInt() int {
	var res int
	if len(d.Price) > 0 {
		p, err := strconv.Atoi(d.Price)
		if err != nil {
			fmt.Println("error price parse", err)
		} else {
			res = p
		}
	}

	return res
}

type QuantityData struct {
	Quantity string `json:"quantity,omitempty" schema:"quantity"`
}

func (d *QuantityData) QuantityInt() int {
	var res int
	if len(d.Quantity) > 0 {
		p, err := strconv.Atoi(d.Quantity)
		if err != nil {
			fmt.Println("error quantity parse", err)
		} else {
			res = p
		}
	}

	return res
}

func (d *QuantityData) CheckQuantityData() bool {
	return len(d.Quantity) > 0
}

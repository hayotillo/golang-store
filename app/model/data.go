package model

import (
	"encoding/json"
	"fmt"
	"store-api/app/misc"
	"strconv"
	"strings"
	"time"
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

type SaleIDData struct {
	SaleID string `json:"sale_id" schema:"sale_id"`
}

func (d *SaleIDData) CheckSaleIDData() bool {
	return len(d.SaleID) == 36
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

type PeriodData struct {
	PeriodBetweenData
	PeriodNameDate
	loaded bool
}

func (d *PeriodData) CheckPeriodData() bool {
	if d.CheckPeriodNameData() || d.CheckPeriodBetweenData() {
		if !d.loaded {
			d.LoadPeriod()
		}
		return true
	}
	return false
}

func (d *PeriodData) PeriodWhere(table string) string {
	if len(d.PeriodStart) < 20 {
		d.PeriodStart = misc.DateParseToFull(d.PeriodStart)
	}
	if len(d.PeriodEnd) != 23 {
		c := misc.CurrentTimeFull()
		d.PeriodEnd = fmt.Sprintf("%s %s", d.PeriodEnd[:10], c[11:])
	}
	if len(d.PeriodStart) == 0 {
		return ""
	}

	return fmt.Sprintf("%[1]s.updated_at::timestamp >= '%[2]s' AND %[1]s.updated_at::timestamp <= '%[3]s'",
		table, d.PeriodStart, d.PeriodEnd)
}

func (d *PeriodData) LoadPeriod() {
	if d.CheckPeriodBetweenData() {
		d.loaded = d.CheckPeriodStartData() && d.CheckPeriodEndData()
		return
	}
	if d.Period == "day" {
		d.PeriodEnd = misc.CurrentTime()
		d.PeriodStart = d.PeriodEnd
		return
	}
	format := "2006-01-02 15:04:05"
	start := ""
	startAt := time.Now()
	d.PeriodEnd = misc.CurrentTimeFull()
	if d.CheckPeriodStartData() {
		if len(d.PeriodStart) == 10 {
			d.PeriodStart = fmt.Sprintf("%s 00:00:00", d.PeriodStart)
		}
		s, err := time.Parse(format, d.PeriodStart)
		if err == nil {
			startAt = s
			d.PeriodEnd = s.Format(format)
		}
	}
	switch d.Period {
	case "day":
		start = startAt.AddDate(0, 0, -1).Format(format)
		break
	case "week":
		start = startAt.AddDate(0, 0, -7).Format(format)
		break
	case "month":
		start = startAt.AddDate(0, -1, 0).Format(format)
		break
	case "quarterly":
		start = startAt.AddDate(0, -3, 0).Format(format)
		break
	case "year":
		start = startAt.AddDate(-1, 0, 0).Format(format)
		break
	}

	d.PeriodStart = start
	d.loaded = true
}

type PeriodEndData struct {
	PeriodEnd string `json:"period_end,omitempty" schema:"period_end"`
}

func (d *PeriodEndData) CheckPeriodEndData() bool {
	if len(d.PeriodEnd) == 10 {
		date, err := time.Parse("2006-01-02", d.PeriodEnd)
		if err != nil {
			fmt.Println("period end date error", err)
			return false
		}
		d.PeriodEnd = date.Format("2006-01-02")
		return true
	}
	return false
}

type PeriodBetweenData struct {
	PeriodStartData
	PeriodEndData
}

func (d *PeriodBetweenData) BetweenTimes() (string, string) {
	bStart := ""
	bEnd := ""
	format := "2006-01-02"
	start, err := time.Parse(format, d.PeriodStart)
	if err == nil {
		bStart = start.Format(format)
	}
	end, err := time.Parse(format, d.PeriodEnd)
	if err == nil {
		bEnd = end.Format(format)
	}

	return bStart, bEnd
}

func (d *PeriodBetweenData) CheckPeriodBetweenData() bool {
	return len(d.PeriodStart) > 9 && len(d.PeriodEnd) > 9
}

type PeriodStartData struct {
	PeriodStart string `json:"period_start,omitempty" schema:"period_start"`
}

func (d *PeriodStartData) CheckPeriodStartData() bool {
	return len(d.PeriodStart) > 9
}

type PeriodNameDate struct {
	Period string `json:"period,omitempty" schema:"period"`
}

func (d *PeriodNameDate) CheckPeriodNameData() bool {
	switch d.Period {
	case "all", "year", "quarterly", "month", "week", "day":
		return true
	default:
		return false
	}
}

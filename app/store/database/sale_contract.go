package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-pdf/fpdf"
	"os"
	"store-api/app/misc"
	"store-api/app/model"
	"store-api/app/store"
	"strconv"
	"strings"
)

type SaleContract struct {
	database *Database
}

func (s *SaleContract) History(f model.SaleHistoryFilter) (*model.SaleHistories, error) {
	if !f.CheckPeriodData() {
		return nil, store.ErrRequiredDataNotFount
	}
	res := &model.SaleHistories{}

	w := ""
	query := fmt.Sprintf(`SELECT
    	s.id,
    	u.id,
    	u.full_name,
    	u.phone,
    	u.status,
    	s.created_at,
    	s.updated_at,
    	sum(sp.price),
    	sum(sp.quantity)
	FROM sales s
		JOIN users u on u.id = s.user_id
		JOIN sale_products sp on s.id = sp.sale_id
		WHERE %s
		GROUP BY s.id, u.id`, f.PeriodWhere("s"))

	if f.CheckPeriodData() {
		w = fmt.Sprintf("%s AND %s", w, f.PeriodWhere("s"))
	}

	rows, err := s.database.db.Query(query)
	//fmt.Println("list", query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var m model.SaleHistoryItem
		rows.Scan(
			&m.ID,
			&m.User.ID,
			&m.User.FullName,
			&m.User.Phone,
			&m.User.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Price,
			&m.Quantity,
		)

		res.Items = append(res.Items, m)
	}

	return res, nil
}

func (s *SaleContract) CheckFile(f model.SaleOneFilter) (*model.ReportFileData, error) {
	if !f.CheckIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	file := &model.ReportFileData{}
	cnf := misc.GetConfig([]string{"company name"})

	// Create a new PDF document
	pdf := fpdf.New(fpdf.OrientationPortrait, fpdf.UnitPoint, fpdf.PageSizeA4, ".")
	w, _ := pdf.GetPageSize()

	// Add a page to the PDF
	pdf.AddPage()
	// set list title
	pdf.SetFont("Arial", "B", 20)
	pdf.SetDrawColor(219, 219, 219)
	pdf.SetTextColor(57, 49, 47)
	// company name
	pdf.CellFormat(w-50, 50, cnf["company name"], "", 0, "C", false, 0, "")
	pdf.Ln(80)

	// Add rows to the table
	pf := model.SaleProductListFilter{}
	pf.PaginatePer = 1000000
	pf.SelfOnly = "true"
	pf.SaleID = f.ID
	products, err := s.Products(pf)
	if err != nil {
		return nil, err
	}

	pdf.SetFont("Arial", "B", 17)
	pdf.CellFormat(300, 30, "Maxsulot", "1", 0, "C", false, 0, "")
	pdf.CellFormat(120, 30, "Narx", "1", 0, "C", false, 0, "")
	pdf.CellFormat(120, 30, "Summa", "1", 0, "C", false, 0, "")

	// row title style
	pdf.SetFont("Arial", "B", 12)
	priceSum := 0
	quantitySum := 0
	for _, item := range products.Items {
		saleProduct := item.(model.SaleProduct)
		pdf.Ln(-1)
		price := saleProduct.PriceInt()
		quantity := saleProduct.QuantityInt()
		pdf.CellFormat(300, 25, saleProduct.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(120, 25, fmt.Sprintf("%vx%v", price, quantity), "1", 0, "R", false, 0, "")
		pdf.CellFormat(120, 25, strconv.Itoa(price*quantity), "1", 0, "R", false, 0, "")
		priceSum += price
		quantitySum += quantity
	}
	pdf.Ln(50)
	pdf.CellFormat(300, 25, "Jami:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(120, 25, strconv.Itoa(priceSum*quantitySum), "1", 0, "R", false, 0, "")
	pdf.CellFormat(120, 25, strconv.Itoa(quantitySum), "1", 0, "R", false, 0, "")
	// add time
	pdf.Ln(50)
	pdf.CellFormat(540, 30, fmt.Sprintf("Sana: %s", misc.CurrentTimeFull()[:19]), "", 0, "C", false, 0, "")

	// Output the PDF to a file
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		fmt.Println("pdf out error:", err)
		os.Exit(1)
	}
	file.Name = strings.Replace(misc.CurrentTimeFull()[:19], " ", "-", -1)
	file.File = buf.Bytes()
	defer pdf.Close()

	return file, nil
}

func (s *SaleContract) Get(f model.SaleOneFilter) (*model.Sale, error) {
	if !f.CheckIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.Sale{}

	query := `SELECT
    	s.id,
    	s.description,
    	s.archive,
    	u.id,
    	u.full_name,
    	u.phone,
    	u.status,
    	s.created_at,
    	s.updated_at,
    	sum(sp.price)
	FROM sales s
		JOIN users u on u.id = s.user_id
		JOIN sale_products sp on s.id = sp.sale_id
	WHERE s.id=$1 GROUP BY s.id, u.id`

	row := s.database.db.QueryRow(query, f.ID)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var description sql.NullString
	row.Scan(
		&m.ID,
		&description,
		&m.Archive,
		&m.User.ID,
		&m.User.FullName,
		&m.User.Phone,
		&m.User.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.Price,
	)
	m.Description = description.String

	return m, nil
}

func (s *SaleContract) One(f model.SaleOneFilter) (*model.SaleData, error) {
	if !f.CheckIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.SaleData{}

	query := `SELECT
    	s.id,
    	s.user_id,
    	s.description,
    	s.archive,
    	s.created_at,
    	s.updated_at
	FROM sales s WHERE s.id=$1`

	row := s.database.db.QueryRow(query, f.ID)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var description sql.NullString
	row.Scan(
		&m.ID,
		&m.UserID,
		&description,
		&m.Archive,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	m.Description = description.String

	return m, nil
}

func (s *SaleContract) List(f model.SaleListFilter) (*model.ListData, error) {
	res := &model.ListData{}

	w := ""
	p := []any{f.Per(), f.Offset()}
	query := `SELECT
    	s.id,
    	s.description,
    	s.archive,
    	u.id,
    	u.full_name,
    	u.phone,
    	u.status,
    	s.created_at,
    	s.updated_at,
    	sum(sp.price),
       	count(s.*) over()
	FROM sales s
		JOIN users u on u.id = s.user_id
		JOIN sale_products sp on s.id = sp.sale_id`

	if f.CheckSearchData() {
		w = fmt.Sprintf("%s AND %s", w, f.SearchLikes([]string{"c.id", "c.full_name", "c.birth", "c.phone"}))
	}

	if f.CheckPeriodData() {
		w = fmt.Sprintf("%s AND %s", w, f.PeriodWhere("s"))
	}

	if len(w) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.TrimPrefix(w, " AND "))
	}
	query = fmt.Sprintf("%s GROUP BY s.id, u.id LIMIT $1 OFFSET $2", query)
	rows, err := s.database.db.Query(query, p...)
	//fmt.Println("list", query)
	if err != nil {
		return nil, err
	}
	count := 0
	for rows.Next() {
		var m model.Sale
		var description sql.NullString
		rows.Scan(
			&m.ID,
			&description,
			&m.Archive,
			&m.User.ID,
			&m.User.FullName,
			&m.User.Phone,
			&m.User.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Price,
			&count,
		)
		m.Description = description.String

		res.Items = append(res.Items, m)
	}
	res.Paginate = f.Paginate(count)

	return res, nil
}

func (s *SaleContract) Products(f model.SaleProductListFilter) (*model.ListData, error) {
	res := &model.ListData{}

	w := ""
	p := []any{f.Per(), f.Offset()}
	groupBy := []string{
		"id",
		"name",
	}
	query := `SELECT p.id, p.name, 0, 0, count(*) over() FROM products p`
	if f.CheckSaleIDData() {
		p = append(p, f.SaleID)
		if f.IsSelfOnly() {
			groupBy = append(groupBy, "sp.price", "sp.quantity")
			query = `SELECT p.id, p.name, sp.price, sp.quantity, count(*) over()
			FROM products p
				JOIN sale_products sp ON p.id = sp.product_id
				WHERE sp.sale_id=$3`
		} else {
			query = `SELECT id, name, MAX(price), MAX(quantity), count(*) over()
			FROM (
				 SELECT p.id, p.name, sp.price, sp.quantity
				 FROM products p
					  JOIN sale_products sp ON p.id = sp.product_id
				 WHERE sp.sale_id = $3
				 UNION ALL
				 SELECT p.id, p.name, 0, 0
				 FROM products p
			) AS combined_data`
		}
	}

	if f.CheckSearchData() {
		w = fmt.Sprintf("%s AND %s", w, f.SearchLikes([]string{"id", "name"}))
	}

	o := f.OrdersWhere([]string{"length(name)", "name"}, "1:asc,2:asc")

	if len(w) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.TrimPrefix(w, " AND "))
	}
	query = fmt.Sprintf("%s GROUP BY %s ORDER BY%s LIMIT $1 OFFSET $2",
		query, strings.TrimSuffix(strings.Join(groupBy, ", "), ", "), o)
	rows, err := s.database.db.Query(query, p...)
	//fmt.Println("list", query)
	if err != nil {
		return nil, err
	}
	count := 0
	for rows.Next() {
		var m model.SaleProduct
		rows.Scan(
			&m.ID,
			&m.Name,
			&m.Price,
			&m.Quantity,
			&count,
		)
		res.Items = append(res.Items, m)
	}
	res.Paginate = f.Paginate(count)

	return res, nil
}

func (s *SaleContract) Save(m *model.SaleData) error {
	if !m.CheckInsertData() && !m.CheckUpdateData() {
		return store.ErrRequiredDataNotFount
	}
	isNew := !m.CheckIDData()
	if isNew {
		m.ID = misc.NewUUID()
	}
	p := []any{m.ID}
	fields := "id,"
	set := ""
	// user
	if m.CheckUserIDData() {
		p = append(p, m.UserID)
		fields = fmt.Sprintf("%s user_id,", fields)
		set = fmt.Sprintf("%s user_id=$%d,", set, len(p))
	}
	// description
	if m.CheckDescriptionData() {
		p = append(p, m.Description)
		fields = fmt.Sprintf("%s description,", fields)
		set = fmt.Sprintf("%s description=$%d,", set, len(p))
	}
	// archive
	if m.CheckArchiveData() {
		p = append(p, m.IsArchive())
		fields = fmt.Sprintf("%s archive,", fields)
		set = fmt.Sprintf("%s archive=$%d,", set, len(p))
	}

	query := fmt.Sprintf(`INSERT INTO sales (%s) VALUES (%s)
	ON CONFLICT (id) DO UPDATE SET%s`,
		strings.TrimSuffix(fields, ","),
		misc.SQLPlaceHolder(len(p), 1),
		strings.TrimSuffix(set, ","))
	_, err := s.database.db.Exec(query, p...)
	if err != nil {
		return err
	}
	one, err := s.One(model.SaleOneFilter{IDData: m.IDData})
	m.UserID = one.UserID
	m.Description = one.Description
	m.Archive = one.Archive
	m.TimestampData = one.TimestampData

	if m.CheckSaleOrdersData() || m.CheckClearOrdersData() {
		query = "DELETE FROM sale_products WHERE sale_id=$1"
		_, err = s.database.db.Exec(query, m.ID)
		if err != nil {
			return err
		}
		if m.CheckSaleOrdersData() {
			p = []any{}
			values := ""
			for _, o := range m.SaleOrders {
				l := len(p)
				p = append(p, m.ID, o.ProductID, o.PriceInt(), o.QuantityInt())
				values = fmt.Sprintf("%s(%s),", values, misc.SQLPlaceHolder(4, 1+l))
			}

			query = fmt.Sprintf(`INSERT INTO sale_products
    		(sale_id, product_id, price, quantity) VALUES %s`,
				strings.TrimSuffix(values, ","))
			_, err = s.database.db.Exec(query, p...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SaleContract) Delete(f model.SaleDeleteFilter) error {
	if !f.HasIDS() {
		return store.ErrRequiredDataNotFount
	}
	var p []any
	for _, id := range f.IDS() {
		p = append(p, id)
	}
	query := fmt.Sprintf("DELETE FROM products WHERE id in (%s)",
		misc.SQLPlaceHolder(len(p), 1))
	res, err := s.database.db.Exec(query, p...)
	if err != nil {
		return err
	}
	if res != nil {
		return nil
	}

	return nil
}

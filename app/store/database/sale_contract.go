package database

import (
	"database/sql"
	"fmt"
	"store-api/app/misc"
	"store-api/app/model"
	"store-api/app/store"
	"strings"
)

type SaleContract struct {
	database *Database
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
	if !f.CheckSaleIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	res := &model.ListData{}

	w := ""
	p := []any{f.Per(), f.Offset(), f.SaleID}
	query := `SELECT id, name, MAX(price), MAX(quantity), count(*) over()
		FROM (
			 SELECT id, name, price, quantity
			 FROM products p
				  JOIN sale_products sp ON p.id = sp.product_id
			 WHERE sp.sale_id = $3
			 UNION ALL
			 SELECT id, name, 0, 0
			 FROM products
		) AS combined_data`

	if f.CheckSearchData() {
		w = fmt.Sprintf("%s AND %s", w, f.SearchLikes([]string{"id", "name"}))
	}

	o := f.OrdersWhere([]string{"length(name)", "name"}, "1:asc,2:asc")

	if len(w) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.TrimPrefix(w, " AND "))
	}
	query = fmt.Sprintf("%s GROUP BY id, name ORDER BY%s LIMIT $1 OFFSET $2", query, o)
	rows, err := s.database.db.Query(query, p...)
	//fmt.Println("list", query)
	if err != nil {
		return nil, err
	}
	count := 0
	for rows.Next() {
		var m model.SaleProduct
		rows.Scan(
			&m.ProductID,
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

	if m.CheckSaleOrdersData() {
		if !isNew {
			query = "DELETE FROM sale_products WHERE sale_id=$1"
			_, err = s.database.db.Exec(query, m.ID)
			if err != nil {
				return err
			}
		}

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

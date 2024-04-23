package database

import (
	"database/sql"
	"fmt"
	"store-api/app/misc"
	"store-api/app/model"
	"store-api/app/store"
	"strings"
)

type ProductContract struct {
	database *Database
}

func (s *ProductContract) One(f model.ProductOneFilter) (*model.Product, error) {
	if !f.CheckIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.Product{}

	query := `SELECT
    	p.id,
    	p.name,
		sum(pi.quantity),
		sum(sp.price)
	FROM products p
	    LEFT JOIN product_incoming pi on p.id = pi.product_id
		LEFT JOIN sale_products sp on p.id = sp.product_id
	WHERE p.id=$1
	GROUP BY p.id`

	row := s.database.db.QueryRow(query, f.ID)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	var inSum, slSum sql.NullString
	row.Scan(&m.ID, &m.Name, &inSum, &slSum)
	m.IncomingSum = inSum.String
	m.SaleSum = slSum.String

	return m, nil
}

func (s *ProductContract) List(f model.ProductListFilter) (*model.ListData, error) {
	res := &model.ListData{}

	w := ""
	p := []any{f.Per(), f.Offset()}
	query := `SELECT
		p.id,
		p.name,
		sum(pi.quantity),
		sum(sp.price),
		count(p.*) over()
	FROM products p
		LEFT JOIN product_incoming pi on p.id = pi.product_id AND pi.archive = false
		LEFT JOIN sale_products sp on p.id = sp.product_id
	GROUP BY p.id`

	if f.CheckSearchData() {
		w = fmt.Sprintf("%s AND %s", w, f.SearchLikes([]string{"c.id", "c.full_name", "c.birth", "c.phone"}))
	}

	if len(w) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.TrimPrefix(w, " AND "))
	}
	o := f.OrdersWhere([]string{"length(p.name)", "p.name"}, "1:asc,2:asc")
	query = fmt.Sprintf("%s ORDER BY%s LIMIT $1 OFFSET $2", query, o)
	rows, err := s.database.db.Query(query, p...)
	if err != nil {
		return nil, err
	}
	count := 0
	for rows.Next() {
		var m model.Product
		var inSum, slSum sql.NullString
		rows.Scan(
			&m.ID,
			&m.Name,
			&inSum,
			&slSum,
			&count,
		)
		m.IncomingSum = inSum.String
		m.SaleSum = slSum.String
		res.Items = append(res.Items, m)
	}
	res.Paginate = f.Paginate(count)

	return res, nil
}

func (s *ProductContract) Save(m *model.Product) error {
	if !m.CheckInsertData() && !m.CheckUpdateData() {
		return store.ErrRequiredDataNotFount
	}
	if !m.CheckIDData() {
		m.ID = misc.NewUUID()
	}
	p := []any{m.ID}
	fields := "id,"
	set := ""
	// full name
	if m.CheckNameData() {
		p = append(p, m.Name)
		fields = fmt.Sprintf("%s name,", fields)
		set = fmt.Sprintf("%s name=$%d,", set, len(p))
	}

	query := fmt.Sprintf(`INSERT INTO products (%s) VALUES (%s)
	ON CONFLICT (id) DO UPDATE SET%s`,
		strings.TrimSuffix(fields, ","),
		misc.SQLPlaceHolder(len(p), 1),
		strings.TrimSuffix(set, ","))
	_, err := s.database.db.Exec(query, p...)
	if err != nil {
		if err.Error() == misc.SqlConstraintErrorStr("products_unique") {
			return store.ErrRecordDuplicate
		}
		return err
	}
	one, err := s.One(model.ProductOneFilter{IDData: m.IDData})
	m.Name = one.Name

	return nil
}

func (s *ProductContract) Delete(f model.ProductDeleteFilter) error {
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

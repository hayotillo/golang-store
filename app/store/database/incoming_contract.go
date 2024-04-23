package database

import (
	"database/sql"
	"fmt"
	"store-api/app/misc"
	"store-api/app/model"
	"store-api/app/store"
	"strings"
)

type IncomingContract struct {
	database *Database
}

func (s *IncomingContract) Get(f model.IncomingOneFilter) (*model.Incoming, error) {
	if !f.CheckIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.Incoming{}

	query := `SELECT
    	pi.id,
    	pi.user_id,
    	pi.product_id,
    	pi.price,
    	pi.quantity,
    	pi.description,
    	pi.created_at,
    	pi.updated_at,
    	u.id,
    	u.full_name,
    	u.phone,
    	p.id,
    	p.name
	FROM product_incoming pi
		JOIN users u on u.id = pi.user_id
		JOIN products p on p.id = pi.product_id
	WHERE pi.id=$1 GROUP BY pi.id, u.id, p.id`

	row := s.database.db.QueryRow(query, f.ID)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	var description sql.NullString
	row.Scan(
		&m.ID,
		&m.Price,
		&m.Quantity,
		&description,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.User.ID,
		&m.User.FullName,
		&m.User.Phone,
		&m.Product.ID,
		&m.Product.Name,
	)
	m.Description = description.String

	return m, nil
}

func (s *IncomingContract) One(f model.IncomingOneFilter) (*model.IncomingData, error) {
	if !f.CheckIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.IncomingData{}

	query := `SELECT
    	pi.id,
    	pi.user_id,
    	pi.product_id,
    	pi.price,
    	pi.quantity,
    	pi.description,
    	pi.created_at,
    	pi.updated_at
	FROM product_incoming pi
	WHERE pi.id=$1`

	row := s.database.db.QueryRow(query, f.ID)
	err := row.Err()
	if err != nil {
		return nil, err
	}

	var description sql.NullString
	row.Scan(
		&m.ID,
		&m.UserID,
		&m.ProductID,
		&m.Price,
		&m.Quantity,
		&description,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	m.Description = description.String

	return m, nil
}

func (s *IncomingContract) List(f model.IncomingListFilter) (*model.ListData, error) {
	res := &model.ListData{}

	w := ""
	p := []any{f.Per(), f.Offset()}
	query := `SELECT
		pi.id,
    	pi.user_id,
    	pi.product_id,
    	pi.price,
    	pi.quantity,
    	pi.description,
    	pi.created_at,
    	pi.updated_at,
    	u.id,
    	u.full_name,
    	u.phone,
    	p.id,
    	p.name,
    	count(pi.*) over()
	FROM product_incoming pi
		JOIN users u on u.id = pi.user_id
		JOIN products p on p.id = pi.product_id`

	if f.CheckSearchData() {
		w = fmt.Sprintf("%s AND %s", w, f.SearchLikes([]string{"c.id", "c.full_name", "c.birth", "c.phone"}))
	}

	if f.CheckProductIDData() {
		p = append(p, f.ProductID)
		w = fmt.Sprintf("%s AND p.id=$%d", w, len(p))
	}

	if len(w) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.TrimPrefix(w, " AND "))
	}
	o := f.OrdersWhere([]string{"length(p.name)", "p.name"}, "1:asc,2:asc")
	query = fmt.Sprintf("%s GROUP BY pi.id, u.id, p.id ORDER BY%s LIMIT $1 OFFSET $2", query, o)
	rows, err := s.database.db.Query(query, p...)
	if err != nil {
		return nil, err
	}
	count := 0
	for rows.Next() {
		var m model.Incoming
		var description sql.NullString
		rows.Scan(
			&m.ID,
			&m.Price,
			&m.Quantity,
			&description,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.User.ID,
			&m.User.FullName,
			&m.User.Phone,
			&m.Product.ID,
			&m.Product.Name,
			&count,
		)
		m.Description = description.String

		res.Items = append(res.Items, m)
	}
	res.Paginate = f.Paginate(count)

	return res, nil
}

func (s *IncomingContract) Save(m *model.IncomingData) error {
	if !m.CheckInsertData() && !m.CheckUpdateData() {
		return store.ErrRequiredDataNotFount
	}

	if !m.CheckIDData() {
		m.ID = misc.NewUUID()
	}
	p := []any{m.ID}
	fields := "id,"
	set := ""
	// price
	if m.CheckPriceData() {
		p = append(p, m.PriceInt())
		fields = fmt.Sprintf("%s price,", fields)
		set = fmt.Sprintf("%s price=$%d,", set, len(p))
	}
	// quantity
	if m.CheckQuantityData() {
		p = append(p, m.QuantityInt())
		fields = fmt.Sprintf("%s quantity,", fields)
		set = fmt.Sprintf("%s quantity=$%d,", set, len(p))
	}
	// description
	if m.CheckDescriptionData() {
		p = append(p, m.Description)
		fields = fmt.Sprintf("%s description,", fields)
		set = fmt.Sprintf("%s description=$%d,", set, len(p))
	}

	query := fmt.Sprintf(`INSERT INTO product_incoming (%s) VALUES (%s)
		ON CONFLICT (id) DO UPDATE SET%s`,
		strings.TrimSuffix(fields, ","),
		misc.SQLPlaceHolder(len(p), 1),
		strings.TrimSuffix(set, ","))
	_, err := s.database.db.Exec(query, p...)
	if err != nil {
		return err
	}

	return nil
}

func (s *IncomingContract) Delete(f model.IncomingDeleteFilter) error {
	if !f.HasIDS() {
		return store.ErrRequiredDataNotFount
	}
	var p []any
	for _, id := range f.IDS() {
		p = append(p, id)
	}
	query := fmt.Sprintf("DELETE FROM product_incoming WHERE id in (%s)",
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

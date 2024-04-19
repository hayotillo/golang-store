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
    	s.status,
    	s.description,
    	s.archive,
    	p.id,
    	p.name,
    	u.id,
    	u.full_name,
    	u.phone,
    	u.status,
    	s.created_at,
    	s.updated_at
	FROM sales s
		JOIN products p on p.id = s.product_id
		JOIN users u on u.id = s.user_id
	WHERE p.id=$1`

	row := s.database.db.QueryRow(query, f.ID)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var description sql.NullString
	row.Scan(
		&m.ID,
		&m.Status,
		&description,
		&m.Archive,
		&m.Product.ID,
		&m.Product.Name,
		&m.User.ID,
		&m.User.FullName,
		&m.User.Phone,
		&m.User.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
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
    	s.product_id,
    	s.status,
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
		&m.ProductID,
		&m.Status,
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
    	s.status,
    	s.description,
    	s.archive,
    	p.id,
    	p.name,
    	u.id,
    	u.full_name,
    	u.phone,
    	u.status,
    	s.created_at,
    	s.updated_at,
       	count(p.*) over()
	FROM sales s
		JOIN products p on p.id = s.product_id
		JOIN users u on u.id = s.user_id`

	if f.CheckSearchData() {
		w = fmt.Sprintf("%s AND %s", w, f.SearchLikes([]string{"c.id", "c.full_name", "c.birth", "c.phone"}))
	}

	if len(w) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.TrimPrefix(w, " AND "))
	}
	query = fmt.Sprintf("%s LIMIT $1 OFFSET $2", query)
	rows, err := s.database.db.Query(query, p...)
	if err != nil {
		return nil, err
	}
	count := 0
	for rows.Next() {
		var m model.Sale
		var description sql.NullString
		rows.Scan(
			&m.ID,
			&m.Status,
			&description,
			&m.Archive,
			&m.Product.ID,
			&m.Product.Name,
			&m.User.ID,
			&m.User.FullName,
			&m.User.Phone,
			&m.User.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
			&count,
		)
		m.Description = description.String

		res.Items = append(res.Items, m)
	}
	res.Paginate = f.Paginate(count)

	return res, nil
}

func (s *SaleContract) Save(m *model.SaleData) error {
	if !m.CheckInsertData() && !m.CheckUpdateData() {
		return store.ErrRequiredDataNotFount
	}
	if !m.CheckIDData() {
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
	// product
	if m.CheckProductIDData() {
		p = append(p, m.ProductID)
		fields = fmt.Sprintf("%s product_id,", fields)
		set = fmt.Sprintf("%s product_id=$%d,", set, len(p))
	}
	// status
	if m.CheckStatusData() {
		p = append(p, m.Status)
		fields = fmt.Sprintf("%s status,", fields)
		set = fmt.Sprintf("%s status=$%d,", set, len(p))
	}
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
		if err.Error() == misc.SqlConstraintErrorStr("products_unique") {
			return store.ErrRecordDuplicate
		}
		return err
	}
	one, err := s.One(model.SaleOneFilter{IDData: m.IDData})
	m.UserID = one.UserID
	m.ProductID = one.ProductID
	m.Status = one.Status
	m.Description = one.Description
	m.Archive = one.Archive
	m.TimestampData = one.TimestampData
	fmt.Println("save", one.Archive)

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

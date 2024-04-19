package database

import (
	"database/sql"
	"fmt"
	"store-api/app/misc"
	"store-api/app/model"
	"store-api/app/store"
	"strings"
)

type VisitContract struct {
	database *Database
}

func (s *VisitContract) Get(f model.VisitOneFilter) (*model.Visit, error) {
	if !f.CheckIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.Visit{}

	query := `SELECT
    	v.id,
    	v.status,
    	v.diagnosis,
    	v.archive,
    	c.id,
    	c.full_name,
    	c.birth,
    	c.phone,
    	c.created_at,
    	c.updated_at,
    	u.id,
    	u.full_name,
    	u.phone,
    	u.status,
    	v.created_at,
    	v.updated_at
	FROM visits v
		JOIN customers c on c.id = v.customer_id
		JOIN users u on u.id = v.user_id
	WHERE c.id=$1`

	row := s.database.db.QueryRow(query, f.ID)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var diagnosis, phone sql.NullString
	row.Scan(
		&m.ID,
		&m.Status,
		&diagnosis,
		&m.Archive,
		&m.Customer.ID,
		&m.Customer.FullName,
		&m.Customer.Birth,
		&phone,
		&m.Customer.CreatedAt,
		&m.Customer.UpdatedAt,
		&m.User.ID,
		&m.User.FullName,
		&m.User.Phone,
		&m.User.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	m.Diagnosis = diagnosis.String
	m.Customer.Phone = phone.String

	return m, nil
}

func (s *VisitContract) One(f model.VisitOneFilter) (*model.VisitData, error) {
	if !f.CheckIDData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.VisitData{}

	query := `SELECT
    	v.id,
    	v.user_id,
    	v.customer_id,
    	v.status,
    	v.diagnosis,
    	v.archive,
    	v.created_at,
    	v.updated_at
	FROM visits v WHERE v.id=$1`

	row := s.database.db.QueryRow(query, f.ID)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var diagnosis sql.NullString
	row.Scan(
		&m.ID,
		&m.UserID,
		&m.CustomerID,
		&m.Status,
		&diagnosis,
		&m.Archive,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	m.Diagnosis = diagnosis.String

	return m, nil
}

func (s *VisitContract) List(f model.VisitListFilter) (*model.ListData, error) {
	res := &model.ListData{}

	w := ""
	p := []any{f.Per(), f.Offset()}
	query := `SELECT
    	v.id,
    	v.status,
    	v.diagnosis,
    	v.archive,
    	c.id,
    	c.full_name,
    	c.birth,
    	c.phone,
    	c.created_at,
    	c.updated_at,
    	u.id,
    	u.full_name,
    	u.phone,
    	u.status,
    	v.created_at,
    	v.updated_at,
       	count(c.*) over()
	FROM visits v
		JOIN customers c on c.id = v.customer_id
		JOIN users u on u.id = v.user_id`

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
		var m model.Visit
		var diagnosis, phone sql.NullString
		rows.Scan(
			&m.ID,
			&m.Status,
			&diagnosis,
			&m.Archive,
			&m.Customer.ID,
			&m.Customer.FullName,
			&m.Customer.Birth,
			&phone,
			&m.Customer.CreatedAt,
			&m.Customer.UpdatedAt,
			&m.User.ID,
			&m.User.FullName,
			&m.User.Phone,
			&m.User.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
			&count,
		)
		m.Diagnosis = diagnosis.String
		m.Customer.Phone = phone.String

		res.Items = append(res.Items, m)
	}
	res.Paginate = f.Paginate(count)

	return res, nil
}

func (s *VisitContract) Save(m *model.VisitData) error {
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
	// customer
	if m.CheckCustomerIDData() {
		p = append(p, m.CustomerID)
		fields = fmt.Sprintf("%s customer_id,", fields)
		set = fmt.Sprintf("%s customer_id=$%d,", set, len(p))
	}
	// status
	if m.CheckStatusData() {
		p = append(p, m.Status)
		fields = fmt.Sprintf("%s status,", fields)
		set = fmt.Sprintf("%s status=$%d,", set, len(p))
	}
	// diagnosis
	if m.CheckDiagnosisData() {
		p = append(p, m.Diagnosis)
		fields = fmt.Sprintf("%s diagnosis,", fields)
		set = fmt.Sprintf("%s diagnosis=$%d,", set, len(p))
	}
	// archive
	if m.CheckArchiveData() {
		p = append(p, m.IsArchive())
		fields = fmt.Sprintf("%s archive,", fields)
		set = fmt.Sprintf("%s archive=$%d,", set, len(p))
	}

	query := fmt.Sprintf(`INSERT INTO visits (%s) VALUES (%s)
	ON CONFLICT (id) DO UPDATE SET%s`,
		strings.TrimSuffix(fields, ","),
		misc.SQLPlaceHolder(len(p), 1),
		strings.TrimSuffix(set, ","))
	_, err := s.database.db.Exec(query, p...)
	if err != nil {
		if err.Error() == misc.SqlConstraintErrorStr("customers_unique") {
			return store.ErrRecordDuplicate
		}
		return err
	}
	one, err := s.One(model.VisitOneFilter{IDData: m.IDData})
	m.UserID = one.UserID
	m.CustomerID = one.CustomerID
	m.Status = one.Status
	m.Diagnosis = one.Diagnosis
	m.Archive = one.Archive
	m.TimestampData = one.TimestampData
	fmt.Println("save", one.Archive)

	return nil
}

func (s *VisitContract) Delete(f model.VisitDeleteFilter) error {
	if !f.HasIDS() {
		return store.ErrRequiredDataNotFount
	}
	var p []any
	for _, id := range f.IDS() {
		p = append(p, id)
	}
	query := fmt.Sprintf("DELETE FROM customers WHERE id in (%s)",
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

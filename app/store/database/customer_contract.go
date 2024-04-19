package database

import (
	"database/sql"
	"fmt"
	"reception/app/misc"
	"reception/app/model"
	"reception/app/store"
	"strings"
)

type CustomerContract struct {
	database *Database
}

func (s *CustomerContract) Get(f model.CustomerOneFilter) (*model.Customer, error) {
	if !f.CheckIDData() && !f.CheckPhoneData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.Customer{}
	var p []any
	query := `SELECT
    	c.id,
    	c.full_name,
    	c.birth,
    	c.phone
	FROM customers c`
	if f.CheckIDData() {
		p = append(p, f.ID)
		query = fmt.Sprintf("%s WHERE c.id=$%d", query, len(p))
	} else {
		p = append(p, f.Phone)
		query = fmt.Sprintf("%s WHERE c.phone=$%d", query, len(p))
	}

	row := s.database.db.QueryRow(query, p...)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var phone sql.NullString
	row.Scan(
		&m.ID,
		&m.FullName,
		&m.Birth,
		&phone,
	)
	if phone.Valid {
		m.Phone = phone.String
	}

	return m, nil
}

func (s *CustomerContract) One(f model.CustomerOneFilter) (*model.CustomerData, error) {
	if !f.CheckIDData() && !f.CheckPhoneData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.CustomerData{}
	var p []any
	query := `SELECT c.id, c.full_name, c.birth, c.phone FROM customers c`
	if f.CheckIDData() {
		p = append(p, f.ID)
		query = fmt.Sprintf("%s WHERE c.id=$%d", query, len(p))
	} else {
		p = append(p, f.Phone)
		query = fmt.Sprintf("%s WHERE c.phone=$%d", query, len(p))
	}

	row := s.database.db.QueryRow(query, p...)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	var phone sql.NullString
	row.Scan(&m.ID, &m.FullName, &m.Birth, &phone)
	m.Phone = phone.String

	return m, nil
}

func (s *CustomerContract) List(f model.CustomerListFilter) (*model.ListData, error) {
	res := &model.ListData{}

	w := ""
	p := []any{f.Per(), f.Offset()}
	query := `SELECT
    	c.id,
       	c.full_name,
       	c.birth,
       	c.phone,
       	count(c.*) over()
	FROM customers c`

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
		var m model.Customer
		var phone sql.NullString
		rows.Scan(
			&m.ID,
			&m.FullName,
			&m.Birth,
			&phone,
			&count,
		)
		if phone.Valid {
			m.Phone = phone.String
		}
		res.Items = append(res.Items, m)
	}
	res.Paginate = f.Paginate(count)

	return res, nil
}

func (s *CustomerContract) Save(m *model.CustomerData) error {
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
	if m.CheckFullNameData() {
		p = append(p, m.FullName)
		fields = fmt.Sprintf("%s full_name,", fields)
		set = fmt.Sprintf("%s full_name=$%d,", set, len(p))
	}
	// birth
	if m.CheckBirthData() {
		p = append(p, m.Birth)
		fields = fmt.Sprintf("%s birth,", fields)
		set = fmt.Sprintf("%s birth=$%d,", set, len(p))
	}
	// phone
	if m.CheckPhoneData() {
		p = append(p, m.Phone)
		fields = fmt.Sprintf("%s phone,", fields)
		set = fmt.Sprintf("%s phone=$%d,", set, len(p))
	}

	query := fmt.Sprintf(`INSERT INTO customers (%s) VALUES (%s)
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
	one, err := s.One(model.CustomerOneFilter{IDData: m.IDData})
	m.DoctorID = one.DoctorID
	m.FullName = one.FullName
	m.Phone = one.Phone

	return nil
}

func (s *CustomerContract) Delete(f model.CustomerDeleteFilter) error {
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

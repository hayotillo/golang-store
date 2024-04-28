package database

import (
	"database/sql"
	"fmt"
	"store-api/app/misc"
	"store-api/app/model"
	"store-api/app/store"
	"strings"
)

type UserContract struct {
	database *Database
}

func (r *UserContract) List(f model.UserListFilter) (*model.ListData, error) {
	list := &model.ListData{}
	//
	p := []any{f.Per(), f.Offset()}
	query := `SELECT
    	u.id,
    	u.phone,
    	u.full_name,
    	u.status,
    	u.created_at,
    	u.updated_at,
    	COUNT(*) OVER()
	FROM users u`

	w := ""
	if f.CheckSearchData() {
		w = f.SearchLikes([]string{"u.phone", "u.full_name", "u.status"})
	}

	orders := f.OrdersWhere([]string{"u.full_name"}, "1:asc")

	if len(w) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.TrimPrefix(w, " AND "))
	}

	if len(orders) > 0 {
		query = fmt.Sprintf("%s GROUP BY u.id ORDER BY%s", query, orders)
	}

	query = fmt.Sprintf("%s LIMIT $1 OFFSET $2", query)

	//fmt.Println("list", query, p)
	rows, err := r.database.db.Query(query, p...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	count := 0
	for rows.Next() {
		var u model.User
		err = rows.Scan(
			&u.ID,
			&u.Phone,
			&u.FullName,
			&u.Status,
			&u.CreatedAt,
			&u.UpdatedAt,
			&count,
		)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, u)
	}

	list.Paginate = f.Paginate(count)

	return list, nil
}

func (r *UserContract) Save(d model.User) (*model.User, error) {
	if !d.CheckInsertData() && !d.CheckUpdateData() {
		return nil, store.ErrRequiredDataNotFount
	}
	m := &model.User{}

	if !d.CheckIDData() {
		d.ID = misc.NewUUID()
	}

	fields := "id"
	set := ""
	p := []any{d.ID}

	if d.CheckPasswordData() {
		p = append(p, d.PasswordHash())
		fields = fmt.Sprintf("%s, encrypt_password", fields)
		set = fmt.Sprintf("%s, encrypt_password=$%d", set, len(p))
	}
	if d.CheckFullNameData() {
		p = append(p, d.FullName)
		fields = fmt.Sprintf("%s, full_name", fields)
		set = fmt.Sprintf("%s, full_name=$%d", set, len(p))
	}
	if d.CheckPhoneData() {
		p = append(p, d.Phone)
		fields = fmt.Sprintf("%s, phone", fields)
		set = fmt.Sprintf("%s, phone=$%d", set, len(p))
	}
	if d.CheckUserStatusData() {
		p = append(p, d.Status)
		fields = fmt.Sprintf("%s, status", fields)
		set = fmt.Sprintf("%s, status=$%d", set, len(p))
	}
	query := fmt.Sprintf(
		`INSERT INTO users (%s) values (%s) ON CONFLICT (id) DO UPDATE SET %s`,
		strings.TrimPrefix(fields, ", "),
		misc.SQLPlaceHolder(len(p), 1),
		strings.TrimPrefix(set, ", "),
	)

	_, err := r.database.db.Exec(query, p...)
	//fmt.Println("user save", err, query, p)
	if err != nil {
		if err.Error() == misc.SqlConstraintErrorStr("users_phone_key") {
			return nil, store.ErrRecordAlreadyExists
		}
		return m, err
	}

	one, err := r.One(model.UserOneFilter{IDData: d.IDData})
	if err != nil {
		return m, err
	}
	m.ID = one.ID
	m.FullName = one.FullName
	m.Phone = one.Phone
	m.StatusData = one.StatusData
	m.Token = one.Token
	m.TimestampData = one.TimestampData

	return m, nil
}

func (r *UserContract) Delete(f model.UserDeleteFilter) (bool, error) {
	if !f.CheckUserIDData() || !f.HasIDS() {
		return false, store.ErrRequiredDataNotFount
	}

	h := ""
	var p []any
	for i, id := range f.IDS() {
		p = append(p, id)
		h = fmt.Sprintf("%s$%d, ", h, i+len(p))
	}

	query := fmt.Sprintf(
		"DELETE FROM users WHERE id IN (%s)",
		strings.TrimSuffix(h, ", "))
	res, err := r.database.db.Exec(query, p...)
	//fmt.Println("delete", err, query, p)
	if err != nil {
		return false, err
	}
	if res != nil {
		return true, nil
	}

	return false, nil
}

func (r *UserContract) SetStatus(d model.UserStatusFilterData) (bool, error) {
	if !d.CheckUserStatusData() || !d.HasIDS() {
		return false, store.ErrRequiredDataNotFount
	}

	h := ""
	p := []any{d.StatusData}
	for i, id := range d.IDS() {
		p = append(p, id)
		fmt.Sprintf("%s$%d, ", h, i+1)
	}

	query := fmt.Sprintf(
		"UPDATE users SET status = $1 WHERE id IN (%s)",
		strings.TrimSuffix(h, ", "))
	res, err := r.database.db.Exec(query, p...)
	if err != nil {
		return false, err
	}

	return res != nil, nil
}

func (r *UserContract) One(f model.UserOneFilter) (*model.User, error) {
	m := &model.User{}
	if !f.CheckIDData() && !f.CheckPhoneData() {
		return m, store.ErrRequiredDataNotFount
	}

	var p []any
	query := `SELECT id,
	   	phone,
	   	full_name,
	   	status,
	   	token,
	   	encrypt_password,
	   	created_at,
	   	updated_at
	FROM users`
	if f.CheckIDData() {
		p = append(p, f.ID)
		query = fmt.Sprintf("%s WHERE id=$%d", query, len(p))
	} else {
		p = append(p, f.Phone)
		query = fmt.Sprintf("%s WHERE phone=$%d", query, len(p))
	}

	if f.CheckUserIDData() && f.UserID != f.ID {
		p = append(p, f.UserID)
		query = fmt.Sprintf("%s AND parent=$%d", query, len(p))
	}
	var token, encryptPassword sql.NullString
	err := r.database.db.QueryRow(query, p...).Scan(
		&m.ID,
		&m.Phone,
		&m.FullName,
		&m.Status,
		&token,
		&encryptPassword,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	//fmt.Println("one", query, err)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	m.Token = token.String
	m.EncryptPassword = encryptPassword.String

	return m, nil
}

func (r *UserContract) GetByPhone(d model.PhoneData) (*model.User, error) {
	m := &model.User{}
	if d.CheckPhoneData() {
		query := `SELECT id,
			   phone,
			   full_name,
			   status,
			   token,
			   encrypt_password,
			   created_at,
			   updated_at
			FROM users
			WHERE phone = $1`
		err := r.database.db.QueryRow(query, d.Phone).Scan(
			&m.ID,
			&m.Phone,
			&m.FullName,
			&m.Status,
			&m.Token,
			&m.EncryptPassword,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		//fmt.Println("by phone", err, d.Phone, query)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}

	return m, nil
}

func (r *UserContract) SetToken(d model.IDData, t string) bool {
	if !d.CheckIDData() {
		return false
	}
	res, err := r.database.db.Exec("UPDATE users SET token=$1 WHERE id=$2", t, d.ID)
	if err != nil {
		return false
	}
	return res != nil
}

package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	IDData
	TimestampData
	FullNameData
	PhoneData
	StatusData
	Token           string `json:"token,omitempty"`
	Password        string `json:"-" schema:"password"`
	EncryptPassword string `json:"-" schema:"-"`
}

type UserStatusFilterData struct {
	IDSData
	StatusData
}

type UserLoginData struct {
	PhoneData
	Password string `schema:"password"`
}

type UserListFilter struct {
	PaginateData
	SearchData
	OrdersData
	IDSData
}

type UserOneFilter struct {
	IDData
	UserIDData
	PhoneData
}

type UserDeleteFilter struct {
	UserIDData
	IDSData
}

func (d *User) ToPublic() User {
	return User{
		IDData:        d.IDData,
		PhoneData:     d.PhoneData,
		FullNameData:  d.FullNameData,
		TimestampData: d.TimestampData,
	}
}

func (d *StatusData) CheckUserStatusData() bool {
	res := false

	for _, s := range []string{"admin", "moderator", "company", "user", "manager"} {
		if s == d.Status {
			return true
		}
	}

	return res
}

func (d *User) IsAdmin() bool {
	return d.Status == "admin"
}
func (d *User) PasswordHash() string {
	if len(d.Password) > 0 {
		enc, err := encryptString(d.Password)
		if err != nil {
			return ""
		}
		return enc
	}
	return ""
}

func (d *User) CheckPhoneData() bool {
	length := len(d.Phone)
	return length > 3 && length < 14
}

func (d *User) CheckPasswordData() bool {
	length := len(d.Password)
	return length > 3 && length < 14
}

func (d *User) CheckFullNameData() bool {
	length := len(d.FullName)
	return length > 3 && length < 50
}

func (d *User) CheckUserStatusData() bool {
	for _, s := range []string{"admin", "doctor", "user"} {
		if s == d.Status {
			return true
		}
	}

	return false
}

func (d *User) CheckRegisterData() bool {
	return d.CheckPhoneData() && d.CheckPasswordData()
}

func (d *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(d.EncryptPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

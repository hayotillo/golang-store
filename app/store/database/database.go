package database

import (
	"database/sql"
	"store-api/app/store"
)

type Database struct {
	db  *sql.DB
	per int
	// user
	userContract *UserContract
	// visit
	visitContract *VisitContract
	// customer
	customerContract *CustomerContract
}

func New(db *sql.DB, per int) *Database {
	return &Database{db: db, per: per}
}

// user

func (d *Database) User() store.UserInterface {
	if d.userContract == nil {
		d.userContract = &UserContract{database: d}
	}
	return d.userContract
}

// visit

func (d *Database) Visit() store.VisitInterface {
	if d.visitContract == nil {
		d.visitContract = &VisitContract{database: d}
	}
	return d.visitContract
}

// customer

func (d *Database) Customer() store.CustomerInterface {
	if d.customerContract == nil {
		d.customerContract = &CustomerContract{database: d}
	}
	return d.customerContract
}

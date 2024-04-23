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
	// sale
	saleContract *SaleContract
	// product
	productContract *ProductContract
	// incoming
	incomingContract *IncomingContract
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

// sale

func (d *Database) Sale() store.SaleInterface {
	if d.saleContract == nil {
		d.saleContract = &SaleContract{database: d}
	}
	return d.saleContract
}

// product

func (d *Database) Product() store.ProductInterface {
	if d.productContract == nil {
		d.productContract = &ProductContract{database: d}
	}
	return d.productContract
}

// incoming

func (d *Database) Incoming() store.IncomingInterface {
	if d.incomingContract == nil {
		d.incomingContract = &IncomingContract{database: d}
	}
	return d.incomingContract
}

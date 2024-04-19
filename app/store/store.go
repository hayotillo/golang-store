package store

type Store interface {
	User() UserInterface

	Sale() SaleInterface

	Product() ProductInterface
}

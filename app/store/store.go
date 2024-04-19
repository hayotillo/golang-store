package store

type Store interface {
	User() UserInterface

	Visit() VisitInterface

	Customer() CustomerInterface
}

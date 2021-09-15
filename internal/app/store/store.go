package store

// Store interface
type Store interface {
	Users() UsersRepository
	Transactions() TransactionsRepository
}

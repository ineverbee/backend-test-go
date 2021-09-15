package sqlstore

import (
	"database/sql"

	"github.com/ineverbee/backend-test-go/internal/app/store"
)

// Store sqlstore structure
type Store struct {
	db                     *sql.DB
	usersRepository        *UsersRepository
	transactionsRepository *TransactionsRepository
}

// New creates new sqlstore
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Users return Users repository
func (s *Store) Users() store.UsersRepository {
	if s.usersRepository != nil {
		return s.usersRepository
	}

	s.usersRepository = &UsersRepository{
		store: s,
	}

	return s.usersRepository
}

// Users return Users repository
func (s *Store) Transactions() store.TransactionsRepository {
	if s.transactionsRepository != nil {
		return s.transactionsRepository
	}

	s.transactionsRepository = &TransactionsRepository{
		store: s,
	}

	return s.transactionsRepository
}

package store

import "github.com/ineverbee/backend-test-go/internal/app/models"

// UsersRepository interface
type UsersRepository interface {
	Create(*models.User) error
	Update(*models.User) error
	FindByID(*models.User) error
}

// TransactionsRepository interface
type TransactionsRepository interface {
	CreateTransaction(*models.Transactions) error
	ListOfTransactions(int64, *[]models.Transactions) error
}

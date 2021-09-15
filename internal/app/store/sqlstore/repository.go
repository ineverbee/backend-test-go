package sqlstore

import (
	"github.com/ineverbee/backend-test-go/internal/app/models"
)

// UsersRepository ...
type UsersRepository struct {
	store *Store
}

// TransactionsRepository ...
type TransactionsRepository struct {
	store *Store
}

// Create creates a new user
func (r *UsersRepository) Create(u *models.User) error {
	if err := r.store.db.QueryRow(
		"INSERT INTO users(balance) VALUES($1) RETURNING id",
		u.Balance,
	).Scan(&u.ID); err != nil {
		return err
	}
	return nil
}

// Update updates a user balance
func (r *UsersRepository) Update(u *models.User) error {
	if err := r.store.db.QueryRow(
		"UPDATE users SET balance = $1 WHERE id = $2 RETURNING id",
		u.Balance,
		u.ID,
	).Scan(&u.ID); err != nil {
		return err
	}
	return nil
}

// CreateTransaction creates a new transaction
func (r *TransactionsRepository) CreateTransaction(t *models.Transactions) error {
	if err := r.store.db.QueryRow(
		"INSERT INTO transactions(user_id, balance_change, comment) VALUES($1, $2, $3) RETURNING id",
		t.User_ID,
		t.BalanceChange,
		t.Comment,
	).Scan(&t.ID); err != nil {
		return err
	}
	return nil
}

// ListOfTransactions retrieves list of transactions for one user
func (r *TransactionsRepository) ListOfTransactions(uid int64, t *[]models.Transactions) error {
	rows, err := r.store.db.Query(
		"SELECT * FROM transactions WHERE user_id = $1",
		uid,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var transaction models.Transactions
		err = rows.Scan(
			&transaction.ID,
			&transaction.User_ID,
			&transaction.BalanceChange,
			&transaction.Comment,
			&transaction.CreatedAt,
		)
		if err != nil {
			return err
		}
		*t = append(*t, transaction)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

// FindByID finds User by ID in db
func (r *UsersRepository) FindByID(u *models.User) error {
	if err := r.store.db.QueryRow(
		"SELECT id, balance FROM users WHERE id = $1",
		u.ID,
	).Scan(
		&u.ID,
		&u.Balance,
	); err != nil {
		return err
	}
	return nil
}

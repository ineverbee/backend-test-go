package models

// User ...
type User struct {
	ID        int64  `json:"id"`
	Balance   int    `json:"balance,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// Transactions ...
type Transactions struct {
	ID            int64  `json:"id"`
	User_ID       int64  `json:"user_id,omitempty"`
	BalanceChange int    `json:"balance_change,omitempty"`
	Comment       string `json:"comment,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
}

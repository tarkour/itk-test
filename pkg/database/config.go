package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Wallet struct {
	WalletID      int
	OperationType string
	Amount        int
}

type DBConn struct {
	Pool *pgxpool.Pool
}

func NewDBConn(db *DBConn) *DBConn {
	return &DBConn{
		Pool: db.Pool,
	}
}

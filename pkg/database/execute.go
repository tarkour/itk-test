package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// func (db *DBConn) Execute() (string, error) {

// 	conn, err := db.Pool.Acquire(context.Background())
// 	if err != nil {
// 		return "", fmt.Errorf("acquire connection failed: %w", err)
// 	}
// 	defer conn.Release()

// 	rows, err := conn.Query(context.Background(), "SELECT * FROM wallets;")
// 	if err != nil {
// 		return "", fmt.Errorf("query execution failed: %v", err)
// 	}
// 	defer rows.Close()

// 	columns := rows.FieldDescriptions()
// 	headers := make([]string, len(columns))
// 	for i, col := range columns {
// 		headers[i] = string(col.Name)
// 	}

// 	var results []string
// 	results = append(results, strings.Join(headers, " | "))
// 	results = append(results, strings.Repeat("---", len(headers)))

// 	for rows.Next() {
// 		values, err := rows.Values()
// 		if err != nil {

// 			return "", fmt.Errorf("read data error: %w", err)
// 		}

// 		row := make([]string, len(values))
// 		for i, val := range values {
// 			row[i] = db.formatValue(val)
// 		}
// 		results = append(results, strings.Join(row, " | "))
// 	}

// 	return strings.Join(results, "\n"), nil

// }

// func (db *DBConn) formatValue(val interface{}) string {
// 	switch v := val.(type) {
// 	case pgtype.Numeric:
// 		num := v.Int.Int64()
// 		return fmt.Sprintf("%d", num)
// 	case pgtype.Timestamp:
// 		return v.Time.Format("2006-01-02 15:04:05")
// 	default:
// 		return fmt.Sprintf("%v", v)
// 	}
// }

func (db *DBConn) DBGet(walletid int) (int, error) {

	ctx := context.Background()

	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("acquire connection failed: %w", err)
	}
	defer conn.Release()

	var balance int

	query := fmt.Sprintf("SELECT balance FROM wallets WHERE walletid = %v;", walletid)

	err = conn.QueryRow(ctx, query).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("wallet not found")
		}
		return 0, fmt.Errorf("query failed: %w", err)
	}

	return balance, nil
}

func (db *DBConn) DBPost(wallet *Wallet) error {

	ctx := context.Background()

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentBalance int
	err = tx.QueryRow(ctx,
		"SELECT balance FROM wallets WHERE walletid = $1 FOR UPDATE", wallet.WalletID).Scan(&currentBalance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("wallet not found")
		}
		return fmt.Errorf("select balance failed: %w", err)
	}

	var newBalance int
	switch wallet.OperationType {
	case "DEPOSIT":
		newBalance = currentBalance + wallet.Amount
	case "WITHDRAW":
		if currentBalance < wallet.Amount {
			return fmt.Errorf("insufficient funds: current %v, requested %v",
				currentBalance, wallet.Amount)
		}
		newBalance = currentBalance - wallet.Amount
	default:
		return fmt.Errorf("invalid operation type: %s", wallet.OperationType)
	}

	_, err = tx.Exec(ctx,
		"UPDATE wallets SET balance = $1 WHERE walletid = $2",
		newBalance,
		wallet.WalletID,
	)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil

}

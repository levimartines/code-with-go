package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record and account entries, and update the accounts' balance.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, result.ToAccount, err = transferMoney(ctx, queries, arg.FromAccountID, arg.ToAccountID, arg.Amount)

		return nil
	})
	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	amount int64,
) (fromAccount Account, toAccount Account, err error) {
	if toAccountID > fromAccountID {
		fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: -amount,
			ID:     fromAccountID,
		})
		if err != nil {
			return
		}

		toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: amount,
			ID:     toAccountID,
		})
		if err != nil {
			return
		}
	} else {
		toAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: amount,
			ID:     toAccountID,
		})
		if err != nil {
			return
		}
		fromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: -amount,
			ID:     fromAccountID,
		})
		if err != nil {
			return
		}
	}
	return
}

func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	transactionQueries := New(tx)
	err = fn(transactionQueries)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

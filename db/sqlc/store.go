package db

import (
	"context"
	"database/sql"
	"fmt"
)

// provides all functions to excute db queries and tranctions
type Store struct{
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// execTx excutes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil{
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		
	}
	
	return tx.Commit()
}

// input parameters of the transfer transction
type TransferTxParams struct{
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

// the result of the transfer transction
type TransferTxResult struct{
	Transfer Transfer `json:"transfer"`
	FromAccountID Account `json:"from_account"`
	ToAccountID Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

// Transfer performs a transfer-related operations within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error){
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil{
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		// TODO

		return nil
	})

	return result, err
}
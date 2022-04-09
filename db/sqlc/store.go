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

// executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	// start a new transaction
	tx, err := store.db.BeginTx(ctx, nil)

	// if there's an error, return it
	if err != nil {
		return err
	}

	// get new Queries object from transaction
	q := New(tx)

	// call input function
	err = fn(q)

	// if there is an error, rollback transaction
	if err != nil {
		// if theres a rollback error,
		//return both transaction error and rollback error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		// if rollback is successful, return transaction error
		return err
	}

	// if all transaction operations are successful,
	// commit the transaction and return its error
	return tx.Commit()
}

// contains input parameters of transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"ammount"`
}

// contains result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account_id"`
	ToAccount   Account  `json:"to_account_id"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// performs a money transfer between two accounts
// creates a transfer record, account entries
// and updates both accounts' balance in a single transaction

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	// empty result
	var result TransferTxResult

	// run new database transaction
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// creates new transfer and stores it in result
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		// creates entry for the "from" account
		// amount is negative since money is going out of the account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// creates entry for the "to" account
		// amount is positive since money is going into the account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// TODO: update accounts' balance

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx,
				q,
				arg.FromAccountID,
				-arg.Amount,
				arg.ToAccountID,
				arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx,
				q,
				arg.ToAccountID,
				arg.Amount,
				arg.FromAccountID,
				-arg.Amount,
			)

		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		accountID1,
		amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		accountID2,
		amount2,
	})

	return

}

package db

import (
	"context"
	"fmt"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)
	print("hello")
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">>BEFORE:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Fatalf("THERE WAS AN ERROR: %v", err)
		}

		result := <-results
		if (TransferTxResult{}) == result {
			t.Fatalf("RESULT IS EMPTY")
		}

		// check transfer
		transfer := result.Transfer
		if (Transfer{}) == transfer {
			t.Fatalf("TRANSFER IS EMPTY")
		}
		if account1.ID != transfer.FromAccountID {
			t.Fatalf("DIFFERENT 'FROM ACCOUNT ID'")
		}

		if account2.ID != transfer.ToAccountID {
			t.Fatalf("DIFFERENT 'TO ACCOUNT ID'")
		}

		if amount != transfer.Amount {
			t.Fatalf("DIFFERENT TRANSFER AMOUNT")
		}

		if transfer.ID == 0 {
			t.Fatalf("INVALID TRANSFER ID")
		}

		if transfer.CreatedAt.IsZero() {
			t.Fatalf("INVALID 'CREATED AT' TIMESTAMP")
		}

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		if err != nil {
			t.Fatalf("TRANSFER DOES NOT EXIST")
		}

		// check entries
		fromEntry := result.FromEntry
		if (Entry{}) == fromEntry {
			t.Fatalf("ENTRY IS EMPTY")
		}
		if account1.ID != fromEntry.AccountID {
			t.Fatalf("DIFFERENT ENTRY 'ACCOUNT ID")
		}

		if -amount != fromEntry.Amount {
			t.Fatalf("DIFFERENT ENTRY AMOUNT")
		}

		if fromEntry.ID == 0 {
			t.Fatalf("INVALID ENTRY ID")
		}

		if fromEntry.CreatedAt.IsZero() {
			t.Fatalf("INVALID 'CREATED AT' TIMESTAMP")
		}

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		if err != nil {
			t.Fatalf("THERE WAS AN ERROR")
		}

		toEntry := result.ToEntry
		if (Entry{}) == toEntry {
			t.Fatalf("ENTRY IS EMPTY")
		}

		if account2.ID != toEntry.AccountID {
			t.Fatalf("DIFFERENT ENTRY 'ACCOUNT ID")
		}

		if amount != toEntry.Amount {
			t.Fatalf("DIFFERENT ENTRY AMOUNT")
		}

		if toEntry.ID == 0 {
			t.Fatalf("INVALID ENTRY ID ")
		}

		if toEntry.CreatedAt.IsZero() {
			t.Fatalf("INVALID 'CREATED AT' TIMESTAMP")
		}

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		if err != nil {
			t.Fatalf("THERE WAS AN ERROR: %v", err)
		}

		// TODO: check accounts
		fromAccount := result.FromAccount
		if (Account{}) == fromAccount {
			t.Fatalf("ACCOUNT IS EMPTY")
		}
		if fromAccount.ID != account1.ID {
			t.Fatalf("DIFFERENT ACCOUNT ID")
		}

		toAccount := result.ToAccount
		if (Account{}) == toAccount {
			t.Fatalf("ACCOUNT IS EMPTY")
		}
		if toAccount.ID != account2.ID {
			t.Fatalf("DIFFERENT ACCOUNT ID")
		}
		fmt.Println(">>TX:", fromAccount.Balance, toAccount.Balance)
		// check accounts' balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		if diff1 != diff2 {
			t.Fatalf("DIFFERENT ACCOUNT BALANCES")
		}

		if diff1 <= 0 {
			t.Fatalf("BALANCE DIFFERECE SHOULD BE POSITIVE")
		}

		if diff1%amount != 0 {
			t.Fatalf("BALANCE DIFFERENCE SHOULD BE DIVISIBLE BY AMOUNT")
		}
	}
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	fmt.Println(">> AFTER: ", updatedAccount1.Balance, updatedAccount2.Balance)

	if updatedAccount1.Balance != account1.Balance-int64(n)*amount {
		t.Fatalf("DIFFERENT ACCOUNT BALANCE")
	}

	if updatedAccount2.Balance != account2.Balance+int64(n)*amount {
		t.Fatalf("DIFFERENT ACCOUNT BALANCE")
	}

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	amount := int64(10)
	errs := make(chan error)
	fmt.Println(">> BEFORE: ", account1.Balance, account2.Balance)
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccoutId := account2.ID
		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccoutId = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccoutId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Fatalf("There was an error: %v", err)
		}
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	fmt.Println(">> AFTER: ", updatedAccount1.Balance, updatedAccount2.Balance)

	if updatedAccount1.Balance != account1.Balance {
		t.Fatalf("DIFFERENT ACCOUNT BALANCE")
	}

	if updatedAccount2.Balance != account2.Balance {
		t.Fatalf("DIFFERENT ACCOUNT BALANCE")
	}
}

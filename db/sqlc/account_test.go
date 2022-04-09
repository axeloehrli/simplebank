package db

import (
	"axel/oehrli/db/util"
	"context"
	"database/sql"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}

	if (Account{}) == account {
		t.Fatalf("ACCOUNT IS EMPTY")
	}

	if arg.Owner != account.Owner {
		t.Fatalf("DIFFERENT OWNER")
	}

	if arg.Balance != account.Balance {
		t.Fatalf("DIFFERENT BALANCE")
	}

	if arg.Currency != account.Currency {
		t.Fatalf("DIFFERENT CURRENCY")
	}

	if account.ID == 0 {
		t.Fatalf("INVALID ID")
	}

	if account.CreatedAt.IsZero() {
		t.Fatalf("INVALID CREATED_AT TIMESTAMP")
	}

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}

	if (Account{}) == account2 {
		t.Fatalf("ACCOUNT IS EMPTY")
	}

	if account1.ID != account2.ID {
		t.Fatalf("DIFFERENT ID")
	}

	if account1.Owner != account2.Owner {
		t.Fatalf("DIFFERENT OWNER")
	}

	if account1.Balance != account2.Balance {
		t.Fatalf("DIFFERENT BALANCE")
	}

	if account1.Currency != account2.Currency {
		t.Fatalf("DIFFERENT CURRENCY")
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}

	if (Account{}) == account2 {
		t.Fatalf("ACCOUNT IS EMPTY")
	}

	if account1.ID != account2.ID {
		t.Fatalf("DIFFERENT ID")
	}

	if account1.Owner != account2.Owner {
		t.Fatalf("DIFFERENT OWNER")
	}

	if arg.Balance != account2.Balance {
		t.Fatalf("DIFFERENT BALANCE")
	}

	if account1.Currency != account2.Currency {
		t.Fatalf("DIFFERENT CURRENCY")
	}
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	if err == nil {
		t.Fatalf("Account was not successfully deleted %v", err)
	}

	if err != sql.ErrNoRows {
		t.Fatalf("Unknown error %v", err)
	}

	if (Account{}) != account2 {
		t.Fatalf("Account should be empty")
	}
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := testQueries.ListAccounts(context.Background())

	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	if len(accounts) == 0 {
		t.Fatalf("ACCOUNT LIST IS EMPTY")
	}

	for _, account := range accounts {
		if (Account{}) == account {
			t.Fatalf("ACCOUNT IS EMPTY")
		}
	}
}

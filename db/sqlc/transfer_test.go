package db

import (
	"github.com/axeloehrli/simplebank/db/util"

	"context"
	"testing"
)

func createRandomTransfer(t *testing.T, account1 Account, account2 Account) Transfer {

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	if (Transfer{}) == transfer {
		t.Fatalf("TRANSFER IS EMPTY")
	}

	if arg.FromAccountID != transfer.FromAccountID {
		t.Fatalf("DIFFERENT 'FROM ACCOUNT ID'")
	}

	if arg.ToAccountID != transfer.ToAccountID {
		t.Fatalf("DIFFERENT 'TO ACCOUNT ID'")
	}

	if transfer.ID == 0 {
		t.Fatalf("INVALID ID")
	}

	if transfer.CreatedAt.IsZero() {
		t.Fatalf("INVALED 'CREATED AT'")
	}

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, account1, account2)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	if (Transfer{}) == transfer2 {
		t.Fatalf("TRANSFER IS EMPTY")
	}

	if transfer1.ID != transfer2.ID {
		t.Fatalf("DIFFERENT TRANSFER ID")
	}

	if transfer1.FromAccountID != transfer2.FromAccountID {
		t.Fatalf("DIFFERENT 'TO ACCOUNT ID'")
	}

	if transfer1.ToAccountID != transfer2.ToAccountID {
		t.Fatalf("DIFFERENT 'FROM ACCOUNT ID'")
	}

	if transfer1.Amount != transfer2.Amount {
		t.Fatalf("DIFFERENT TRANSFER AMOUNT")
	}
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
	}

	transfers, err := testQueries.ListTransfers(context.Background())

	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	if len(transfers) == 0 {
		t.Fatalf("TRANSFERS LIST IS EMPTY")
	}

	for _, transfer := range transfers {
		if (Transfer{}) == transfer {
			t.Fatalf("TRANSFER IS EMPTY")
		}
	}
}

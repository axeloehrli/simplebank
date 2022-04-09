package db

import (
	"axel/oehrli/db/util"
	"context"
	"testing"
)

func CreateRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	if (Entry{}) == entry {
		t.Fatalf("ENTRY IS EMPTY")
	}

	if arg.AccountID != entry.AccountID {
		t.Fatalf("DIFFERENT ENTRY ACCOUNT ID")
	}

	if arg.Amount != entry.Amount {
		t.Fatalf("DIFFERENT ENTRY AMOUNT")
	}

	if entry.ID == 0 {
		t.Fatalf("INVALID ENTRY ID")
	}

	if entry.CreatedAt.IsZero() {
		t.Fatalf("INVALID CREATED AT TIMESTAMP")
	}

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	CreateRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := CreateRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	if err != nil {
		t.Fatalf("THERE WAS AN ERROR: %v", err)
	}

	if (Entry{}) == entry2 {
		t.Fatalf("ENTRY IS EMPTY")
	}

	if entry1.ID != entry2.ID {
		t.Fatalf("DIFFERENT ENTRY ID")
	}

	if entry1.AccountID != entry2.AccountID {
		t.Fatalf("DIFFERENT ENTRY ACCOUNT ID")
	}

	if entry1.Amount != entry2.Amount {
		t.Fatalf("DIFFERENT ENTRY AMOUNT")
	}
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		account := createRandomAccount(t)
		CreateRandomEntry(t, account)
	}

	entries, err := testQueries.ListEntries(context.Background())

	if err != nil {
		t.Fatalf("THERE WAS AN ERRO: %v", err)
	}

	if len(entries) == 0 {
		t.Fatalf("LIST OF ENTRIES IS EMPTY")
	}

	for _, entry := range entries {
		if (Entry{}) == entry {
			t.Fatal("ENTRY IS EMPTY")
		}
	}
}

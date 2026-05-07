package repository_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/domain/repository"
	"github.com/jonesrussell/godo/internal/domain/testhelpers"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

func TestNoteRepository_SQLite_CRUD(t *testing.T) {
	t.Parallel()

	log := logger.NewTestLogger(t)
	store, cleanup := testhelpers.NewTempSQLiteUnified(t, log)
	t.Cleanup(cleanup)

	repo := repository.NewNoteRepository(store)
	ctx := context.Background()

	n := model.NewNote("hello sqlite")
	if err := repo.Add(ctx, n); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if n.ID == "" {
		t.Fatal("expected ID after Add")
	}

	got, err := repo.GetByID(ctx, n.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Content != "hello sqlite" || got.Done {
		t.Fatalf("GetByID: %+v", got)
	}

	n.Content = "updated"
	n.Done = true
	if err := repo.Update(ctx, n); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got2, err := repo.GetByID(ctx, n.ID)
	if err != nil {
		t.Fatalf("GetByID after update: %v", err)
	}
	if got2.Content != "updated" || !got2.Done {
		t.Fatalf("after update: %+v", got2)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("List len: %d", len(list))
	}

	if err := repo.Delete(ctx, n.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.GetByID(ctx, n.ID); err == nil {
		t.Fatal("expected error after delete")
	}
}

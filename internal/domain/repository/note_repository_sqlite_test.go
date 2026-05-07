package repository

import (
	"context"
	"testing"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/domain/testfixtures"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/sqlite"
)

func TestNoteRepository_SQLite_CRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := testfixtures.NewTempSQLiteStore(t)
	adapter := sqlite.NewUnifiedAdapter(store)
	repo := NewNoteRepository(adapter)

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
		t.Fatalf("unexpected note: %+v", got)
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
		t.Fatalf("unexpected updated note: %+v", got2)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 note, got %d", len(list))
	}

	if err := repo.Delete(ctx, n.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.GetByID(ctx, n.ID); err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestNoteRepository_SQLite_ListMultiple(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := testfixtures.NewTempSQLiteStore(t)
	adapter := sqlite.NewUnifiedAdapter(store)
	repo := NewNoteRepository(adapter)

	for _, content := range []string{"a", "b"} {
		n := model.NewNote(content)
		if err := repo.Add(ctx, n); err != nil {
			t.Fatalf("Add %q: %v", content, err)
		}
	}
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("len=%d", len(list))
	}
}

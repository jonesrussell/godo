package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domainstorage "github.com/jonesrussell/godo/internal/domain/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/api"
	storageerrors "github.com/jonesrussell/godo/internal/infrastructure/storage/errors"
)

func TestStore_DeleteNote_NotFound(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	log := logger.NewTestLogger(t)
	cfg := domainstorage.APIConfig{
		BaseURL:            srv.URL,
		Timeout:            5,
		RetryCount:         0,
		RetryDelay:         1,
		InsecureSkipVerify: false,
	}

	s, err := api.New(cfg, log)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	err = s.DeleteNote(context.Background(), "missing-id")
	var nf *storageerrors.NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected NotFoundError, got %T %v", err, err)
	}
	if nf.ID != "missing-id" {
		t.Fatalf("unexpected id: %q", nf.ID)
	}
}

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domainstorage "github.com/jonesrussell/godo/internal/domain/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

func TestPatchNoteJSONNoBody(t *testing.T) {
	t.Parallel()
	now := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		_ = json.NewEncoder(w).Encode(APIResponse{
			Data: APINote{
				ID:        "n1",
				Content:   "c",
				Done:      true,
				CreatedAt: now,
				UpdatedAt: now,
			},
		})
	}))
	t.Cleanup(ts.Close)

	log := logger.NewNoopLogger()
	st, err := New(domainstorage.APIConfig{
		BaseURL:               ts.URL,
		Timeout:               5,
		RetryCount:            0,
		RetryDelay:            1,
		TLSInsecureSkipVerify: false,
	}, log)
	if err != nil {
		t.Fatal(err)
	}

	note, err := st.patchNoteJSONNoBody(context.Background(), "abc", "toggle")
	if err != nil {
		t.Fatal(err)
	}
	if note.ID != "n1" || note.Content != "c" || !note.Done {
		t.Fatalf("note=%+v", note)
	}
}

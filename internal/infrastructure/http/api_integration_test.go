package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/jonesrussell/godo/internal/domain/repository"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/domain/testfixtures"
	"github.com/jonesrussell/godo/internal/infrastructure/api"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage/sqlite"
)

func mintTestJWT(secret string) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "test-user",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	s, err := tok.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return s
}

func newTestAPIServer(t *testing.T) (*api.Server, string) {
	t.Helper()
	log := logger.NewNoopLogger()
	st := testfixtures.NewTempSQLiteStore(t)
	adapter := sqlite.NewUnifiedAdapter(st)
	repo := repository.NewNoteRepository(adapter)
	svc := service.NewNoteService(repo, log)
	const secret = "test-secret-for-ci"
	srv := api.NewServer(svc, log, secret)
	return srv, mintTestJWT(secret)
}

func TestAPI_NotesCRUD_JSON(t *testing.T) {
	t.Parallel()
	srv, token := newTestAPIServer(t)
	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	auth := "Bearer " + token
	client := ts.Client()

	// GET /api/v1/notes (empty)
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/notes", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("list status=%d body=%s", resp.StatusCode, b)
	}

	// POST create
	body := strings.NewReader(`{"content":"from test"}`)
	req2, err := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/notes", body)
	if err != nil {
		t.Fatal(err)
	}
	req2.Header.Set("Authorization", auth)
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp2.Body)
		t.Fatalf("create status=%d body=%s", resp2.StatusCode, b)
	}
	var created api.NoteResponse
	if err := json.NewDecoder(resp2.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Content != "from test" {
		t.Fatalf("unexpected created: %+v", created)
	}

	// PUT update
	up := strings.NewReader(`{"content":"updated","done":true}`)
	req3, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/notes/"+created.ID, up)
	if err != nil {
		t.Fatal(err)
	}
	req3.Header.Set("Authorization", auth)
	req3.Header.Set("Content-Type", "application/json")
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp3.Body)
		t.Fatalf("put status=%d body=%s", resp3.StatusCode, b)
	}

	// DELETE
	req4, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/notes/"+created.ID, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	req4.Header.Set("Authorization", auth)
	resp4, err := client.Do(req4)
	if err != nil {
		t.Fatal(err)
	}
	defer resp4.Body.Close()
	if resp4.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp4.Body)
		t.Fatalf("delete status=%d body=%s", resp4.StatusCode, b)
	}
}

func TestAPI_CreateNote_InvalidJSON(t *testing.T) {
	t.Parallel()
	srv, token := newTestAPIServer(t)
	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/notes", bytes.NewBufferString(`{`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFileHandleServesBrotliWithContentEncoding(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "app.js.br"), []byte("fake-brotli"), 0644); err != nil {
		t.Fatal(err)
	}
	initMimeExt()
	_handler = http.FileServer(http.Dir(dir))

	req := httptest.NewRequest(http.MethodGet, "/app.js.br", nil)
	rec := httptest.NewRecorder()
	fileHandle(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "br" {
		t.Fatalf("Content-Encoding = %q, want %q", got, "br")
	}
	if got := rec.Header().Get("Content-Type"); got != "application/javascript" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/javascript")
	}
}

package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const testReport = `{
  "metadata": {
    "generatedAt": "2025-06-15T06:00:00Z",
    "collectorVersion": "test",
    "clusterName": "test-cluster",
    "namespacesScanned": ["default"]
  },
  "images": [
    {
      "image": "nginx:1.27",
      "digest": "sha256:abc123def456",
      "namespace": "default",
      "workload": {"kind": "Deployment", "name": "nginx"},
      "signature": {"signed": true, "verified": true},
      "sbom": {"hasSBOM": true, "format": "spdx"}
    }
  ],
  "helmReleases": [
    {
      "releaseName": "ingress",
      "namespace": "default",
      "chart": "ingress-nginx",
      "version": "4.8.0",
      "appVersion": "1.9.4",
      "status": "deployed"
    }
  ],
  "summary": {
    "totalImages": 1,
    "uniqueImages": 1,
    "signedImages": 1,
    "verifiedImages": 1,
    "imagesWithSBOM": 1,
    "imagesWithUpdates": 0,
    "totalHelmReleases": 1,
    "helmReleasesWithUpdates": 0
  }
}`

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "provenance-20250615-060000.json"), []byte(testReport), 0o644); err != nil {
		t.Fatal(err)
	}
	// Write the latest symlink file too — should be excluded from listing
	if err := os.WriteFile(filepath.Join(dir, "provenance-latest.json"), []byte(testReport), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestHealthz(t *testing.T) {
	srv := NewServer(t.TempDir())
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestListReports(t *testing.T) {
	dir := setupTestDir(t)
	srv := NewServer(dir)

	req := httptest.NewRequest("GET", "/api/reports", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var reports []reportEntry
	if err := json.Unmarshal(w.Body.Bytes(), &reports); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Should have 1 report (provenance-latest.json is excluded)
	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}
	if reports[0].Filename != "provenance-20250615-060000.json" {
		t.Errorf("expected provenance-20250615-060000.json, got %s", reports[0].Filename)
	}
	if reports[0].Summary.TotalImages != 1 {
		t.Errorf("expected totalImages=1, got %d", reports[0].Summary.TotalImages)
	}
	if reports[0].ClusterName != "test-cluster" {
		t.Errorf("expected clusterName=test-cluster, got %s", reports[0].ClusterName)
	}
}

func TestListReports_EmptyDir(t *testing.T) {
	srv := NewServer(t.TempDir())
	req := httptest.NewRequest("GET", "/api/reports", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var reports []reportEntry
	if err := json.Unmarshal(w.Body.Bytes(), &reports); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if reports != nil {
		t.Errorf("expected null/empty, got %d reports", len(reports))
	}
}

func TestGetReport(t *testing.T) {
	dir := setupTestDir(t)
	srv := NewServer(dir)

	req := httptest.NewRequest("GET", "/api/reports/provenance-20250615-060000.json", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestGetReport_NotFound(t *testing.T) {
	srv := NewServer(t.TempDir())
	req := httptest.NewRequest("GET", "/api/reports/nonexistent.json", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetReport_PathTraversal(t *testing.T) {
	srv := NewServer(t.TempDir())
	req := httptest.NewRequest("GET", "/api/reports/..%2f..%2fetc%2fpasswd", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	// The handler checks for ".." and "/" in the filename and rejects them.
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for path traversal, got %d", w.Code)
	}
}

func TestIndex(t *testing.T) {
	srv := NewServer(t.TempDir())
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected text/html, got %s", ct)
	}
	body := w.Body.String()
	if len(body) < 100 {
		t.Error("expected substantial HTML body")
	}
}

func TestIndex_ContainsDetailPanel(t *testing.T) {
	srv := NewServer(t.TempDir())
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	body := w.Body.String()
	for _, expected := range []string{
		"detail-panel",
		"detail-overlay",
		"openDetailIdx",
		"closeDetail",
		"detail-section",
		"lastFilteredImages",
	} {
		if !contains(body, expected) {
			t.Errorf("expected HTML to contain %q", expected)
		}
	}
}

func TestGetReport_ValidJSON(t *testing.T) {
	dir := setupTestDir(t)
	srv := NewServer(dir)

	req := httptest.NewRequest("GET", "/api/reports/provenance-20250615-060000.json", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	var report map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if _, ok := report["metadata"]; !ok {
		t.Error("expected metadata field in report")
	}
	if _, ok := report["images"]; !ok {
		t.Error("expected images field in report")
	}
	if _, ok := report["summary"]; !ok {
		t.Error("expected summary field in report")
	}
}

func TestListReports_SortedNewestFirst(t *testing.T) {
	dir := t.TempDir()
	// Create two reports with different timestamps
	report1 := `{"metadata":{"generatedAt":"2025-06-14T06:00:00Z","collectorVersion":"test","namespacesScanned":["default"]},"images":[],"summary":{"totalImages":0}}`
	report2 := `{"metadata":{"generatedAt":"2025-06-15T06:00:00Z","collectorVersion":"test","namespacesScanned":["default"]},"images":[],"summary":{"totalImages":0}}`

	if err := os.WriteFile(filepath.Join(dir, "provenance-20250614-060000.json"), []byte(report1), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "provenance-20250615-060000.json"), []byte(report2), 0o644); err != nil {
		t.Fatal(err)
	}

	srv := NewServer(dir)
	req := httptest.NewRequest("GET", "/api/reports", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	var reports []reportEntry
	if err := json.Unmarshal(w.Body.Bytes(), &reports); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
	// Newest first
	if reports[0].GeneratedAt < reports[1].GeneratedAt {
		t.Errorf("expected newest first, got %s before %s", reports[0].GeneratedAt, reports[1].GeneratedAt)
	}
}

func TestGetReport_EmptyFilename(t *testing.T) {
	srv := NewServer(t.TempDir())
	req := httptest.NewRequest("GET", "/api/reports/", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty filename, got %d", w.Code)
	}
}

func TestListReports_SkipsMalformedFiles(t *testing.T) {
	dir := t.TempDir()
	// Valid report
	if err := os.WriteFile(filepath.Join(dir, "provenance-20250615-060000.json"), []byte(testReport), 0o644); err != nil {
		t.Fatal(err)
	}
	// Malformed JSON
	if err := os.WriteFile(filepath.Join(dir, "provenance-20250614-060000.json"), []byte("{bad json"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Non-JSON file (should be skipped)
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("not a report"), 0o644); err != nil {
		t.Fatal(err)
	}

	srv := NewServer(dir)
	req := httptest.NewRequest("GET", "/api/reports", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	var reports []reportEntry
	if err := json.Unmarshal(w.Body.Bytes(), &reports); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(reports) != 1 {
		t.Errorf("expected 1 valid report (skipping malformed + non-json), got %d", len(reports))
	}
}

func TestListReports_NonexistentDir(t *testing.T) {
	srv := NewServer("/nonexistent/path")
	req := httptest.NewRequest("GET", "/api/reports", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for nonexistent dir, got %d", w.Code)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

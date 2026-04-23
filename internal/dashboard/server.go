package dashboard

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nebari-dev/provenance-collector/internal/report"
)

// Server serves the provenance dashboard web UI.
type Server struct {
	reportsDir string
	mux        *http.ServeMux
}

// NewServer creates a dashboard HTTP handler that reads reports from reportsDir.
func NewServer(reportsDir string) *Server {
	s := &Server{reportsDir: reportsDir}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/reports", s.handleListReports)
	mux.HandleFunc("/api/reports/", s.handleGetReport)
	mux.HandleFunc("/api/export", s.handleExport)
	mux.HandleFunc("/healthz", s.handleHealthz)
	s.mux = mux
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// reportEntry is a summary of a report file for the listing API.
type reportEntry struct {
	Filename    string               `json:"filename"`
	GeneratedAt string               `json:"generatedAt"`
	Summary     report.ReportSummary `json:"summary"`
	ClusterName string               `json:"clusterName,omitempty"`
}

func (s *Server) handleListReports(w http.ResponseWriter, _ *http.Request) {
	entries, err := os.ReadDir(s.reportsDir)
	if err != nil {
		http.Error(w, "failed to read reports directory", http.StatusInternalServerError)
		return
	}

	var reports []reportEntry
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") || e.Name() == "provenance-latest.json" {
			continue
		}

		r, err := s.loadReport(e.Name())
		if err != nil {
			slog.Warn("skipping malformed report", "file", e.Name(), "error", err)
			continue
		}

		reports = append(reports, reportEntry{
			Filename:    e.Name(),
			GeneratedAt: r.Metadata.GeneratedAt.Format("2006-01-02T15:04:05Z"),
			Summary:     r.Summary,
			ClusterName: r.Metadata.ClusterName,
		})
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].GeneratedAt > reports[j].GeneratedAt
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reports)
}

func (s *Server) handleGetReport(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/reports/")
	if filename == "" || strings.Contains(filename, "/") || strings.Contains(filename, "..") {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	// "latest" resolves to provenance-latest.json
	if filename == "latest" {
		filename = "provenance-latest.json"
	}

	data, err := os.ReadFile(filepath.Join(s.reportsDir, filename))
	if err != nil {
		http.Error(w, "report not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (s *Server) loadReport(filename string) (*report.ProvenanceReport, error) {
	data, err := os.ReadFile(filepath.Join(s.reportsDir, filename))
	if err != nil {
		return nil, err
	}
	var r report.ProvenanceReport
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Server) handleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

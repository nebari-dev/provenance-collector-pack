package report

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func testReport() *ProvenanceReport {
	return &ProvenanceReport{
		Metadata: ReportMetadata{
			GeneratedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			CollectorVersion:  "test",
			ClusterName:       "test-cluster",
			NamespacesScanned: []string{"default"},
		},
		Images: []ImageRecord{
			{
				Image:     "nginx:1.27",
				Digest:    "sha256:abc123",
				Namespace: "default",
				Workload:  WorkloadRef{Kind: "Deployment", Name: "nginx"},
			},
		},
		Summary: ReportSummary{
			TotalImages:  1,
			UniqueImages: 1,
		},
	}
}

func TestPVCWriter(t *testing.T) {
	dir := t.TempDir()
	writer := NewPVCWriter(dir)

	err := writer.Write(context.Background(), testReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that latest file was written
	latestPath := filepath.Join(dir, "provenance-latest.json")
	data, err := os.ReadFile(latestPath)
	if err != nil {
		t.Fatalf("failed to read latest report: %v", err)
	}

	var report ProvenanceReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("failed to unmarshal report: %v", err)
	}

	if report.Metadata.ClusterName != "test-cluster" {
		t.Errorf("expected cluster name test-cluster, got %s", report.Metadata.ClusterName)
	}
	if len(report.Images) != 1 {
		t.Errorf("expected 1 image, got %d", len(report.Images))
	}
}

func TestConfigMapWriter(t *testing.T) {
	client := fake.NewSimpleClientset()
	writer := NewConfigMapWriter(client, "test-report", "default")

	// First write — creates the ConfigMap
	err := writer.Write(context.Background(), testReport())
	if err != nil {
		t.Fatalf("unexpected error on create: %v", err)
	}

	cm, err := client.CoreV1().ConfigMaps("default").Get(context.Background(), "test-report", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get configmap: %v", err)
	}

	if _, ok := cm.Data["report.json"]; !ok {
		t.Error("expected report.json key in configmap data")
	}
	if cm.Labels["app.kubernetes.io/name"] != "provenance-collector" {
		t.Error("expected provenance-collector label")
	}

	// Second write — updates the ConfigMap
	err = writer.Write(context.Background(), testReport())
	if err != nil {
		t.Fatalf("unexpected error on update: %v", err)
	}
}

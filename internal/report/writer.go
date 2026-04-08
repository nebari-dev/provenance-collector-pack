package report

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Writer persists a ProvenanceReport.
type Writer interface {
	Write(ctx context.Context, report *ProvenanceReport) error
}

// PVCWriter writes reports as JSON files to a filesystem path (typically a PVC mount).
type PVCWriter struct {
	basePath string
}

// NewPVCWriter creates a Writer that outputs to the filesystem.
func NewPVCWriter(basePath string) Writer {
	return &PVCWriter{basePath: basePath}
}

func (w *PVCWriter) Write(_ context.Context, report *ProvenanceReport) error {
	if err := os.MkdirAll(w.basePath, 0o755); err != nil {
		return fmt.Errorf("creating report directory: %w", err)
	}

	filename := fmt.Sprintf("provenance-%s.json", time.Now().UTC().Format("20060102-150405"))
	path := filepath.Join(w.basePath, filename)

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling report: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing report to %s: %w", path, err)
	}

	// Also write a "latest" symlink/copy for easy access
	latestPath := filepath.Join(w.basePath, "provenance-latest.json")
	_ = os.Remove(latestPath)
	if err := os.WriteFile(latestPath, data, 0o644); err != nil {
		return fmt.Errorf("writing latest report: %w", err)
	}

	return nil
}

// ConfigMapWriter writes reports as Kubernetes ConfigMaps.
type ConfigMapWriter struct {
	client    kubernetes.Interface
	name      string
	namespace string
}

// NewConfigMapWriter creates a Writer that stores reports in a ConfigMap.
func NewConfigMapWriter(client kubernetes.Interface, name, namespace string) Writer {
	return &ConfigMapWriter{
		client:    client,
		name:      name,
		namespace: namespace,
	}
}

func (w *ConfigMapWriter) Write(ctx context.Context, report *ProvenanceReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling report: %w", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      w.name,
			Namespace: w.namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       "provenance-collector",
				"app.kubernetes.io/managed-by": "provenance-collector",
			},
		},
		Data: map[string]string{
			"report.json": string(data),
			"generated":   report.Metadata.GeneratedAt.Format(time.RFC3339),
		},
	}

	existing, err := w.client.CoreV1().ConfigMaps(w.namespace).Get(ctx, w.name, metav1.GetOptions{})
	if err != nil {
		// Create new
		_, err = w.client.CoreV1().ConfigMaps(w.namespace).Create(ctx, cm, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("creating configmap %s/%s: %w", w.namespace, w.name, err)
		}
		return nil
	}

	// Update existing
	existing.Data = cm.Data
	existing.Labels = cm.Labels
	_, err = w.client.CoreV1().ConfigMaps(w.namespace).Update(ctx, existing, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("updating configmap %s/%s: %w", w.namespace, w.name, err)
	}
	return nil
}

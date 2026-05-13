package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ScanRunner abstracts the side effects of /api/scan so tests can swap in a
// fake without standing up an apiserver.
type ScanRunner interface {
	// HasActiveManualJob reports whether a Job owned by the watched CronJob
	// is currently active (not yet completed or failed).
	HasActiveManualJob(ctx context.Context) (bool, error)
	// StartManualScan creates a new one-shot Job from the watched CronJob's
	// jobTemplate and returns the new Job's name.
	StartManualScan(ctx context.Context) (string, error)
}

// k8sScanRunner is the production ScanRunner: it talks to a real apiserver
// via client-go to inspect and create Jobs from a specific CronJob.
type k8sScanRunner struct {
	client    kubernetes.Interface
	namespace string
	cronJob   string
}

// NewK8sScanRunner builds a ScanRunner backed by the provided clientset.
// Exported so cmd/dashboard can wire it from an in-cluster config.
func NewK8sScanRunner(client kubernetes.Interface, namespace, cronJob string) ScanRunner {
	return &k8sScanRunner{client: client, namespace: namespace, cronJob: cronJob}
}

func (k *k8sScanRunner) HasActiveManualJob(ctx context.Context) (bool, error) {
	jobs, err := k.client.BatchV1().Jobs(k.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, j := range jobs.Items {
		if !ownedBy(j.OwnerReferences, "CronJob", k.cronJob) {
			continue
		}
		// Active counts pods currently running. A freshly created Job may
		// have Active=0 for a heartbeat before pods spin up, so also treat
		// a Job with no terminal state as in-flight.
		if j.Status.Active > 0 {
			return true, nil
		}
		if j.Status.CompletionTime == nil && j.Status.Failed == 0 && j.Status.Succeeded == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (k *k8sScanRunner) StartManualScan(ctx context.Context) (string, error) {
	cj, err := k.client.BatchV1().CronJobs(k.namespace).Get(ctx, k.cronJob, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	name := fmt.Sprintf("manual-%d", time.Now().Unix())
	yes := true
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: k.namespace,
			Labels:    cj.Spec.JobTemplate.Labels,
			Annotations: map[string]string{
				// Same annotation `kubectl create job --from=cronjob/...` sets.
				"cronjob.kubernetes.io/instantiate": "manual",
			},
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         "batch/v1",
				Kind:               "CronJob",
				Name:               cj.Name,
				UID:                cj.UID,
				Controller:         &yes,
				BlockOwnerDeletion: &yes,
			}},
		},
		Spec: cj.Spec.JobTemplate.Spec,
	}
	created, err := k.client.BatchV1().Jobs(k.namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return created.Name, nil
}

func ownedBy(refs []metav1.OwnerReference, kind, name string) bool {
	for _, o := range refs {
		if o.Kind == kind && o.Name == name {
			return true
		}
	}
	return false
}

// scanResponse is the JSON shape returned by a successful POST /api/scan.
type scanResponse struct {
	JobName   string `json:"jobName"`
	Namespace string `json:"namespace"`
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !s.auth.enabled() {
		http.Error(w, "scan endpoint is not configured (oidcIssuer/adminGroups unset)", http.StatusServiceUnavailable)
		return
	}
	id := s.auth.identify(r.Context(), r)
	if !s.auth.canRunScan(id) {
		http.Error(w, "forbidden: caller is not in an admin group", http.StatusForbidden)
		return
	}
	if s.scan == nil {
		http.Error(w, "scan endpoint is not configured (kubernetes client unavailable)", http.StatusServiceUnavailable)
		return
	}

	active, err := s.scan.HasActiveManualJob(r.Context())
	if err != nil {
		http.Error(w, "failed to inspect running jobs: "+err.Error(), http.StatusBadGateway)
		return
	}
	if active {
		http.Error(w, "a scan job is already active for this collector", http.StatusConflict)
		return
	}

	name, err := s.scan.StartManualScan(r.Context())
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "scan request canceled", http.StatusGatewayTimeout)
			return
		}
		http.Error(w, "failed to start scan: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(scanResponse{
		JobName:   name,
		Namespace: namespaceOf(s.scan),
	})
}

// namespaceOf surfaces the namespace from a ScanRunner for the response body
// without leaking the concrete type. Returns "" for runners that don't expose
// one (e.g. test fakes that don't need it).
func namespaceOf(r ScanRunner) string {
	if k, ok := r.(*k8sScanRunner); ok {
		return k.namespace
	}
	if n, ok := r.(interface{ Namespace() string }); ok {
		return n.Namespace()
	}
	return ""
}

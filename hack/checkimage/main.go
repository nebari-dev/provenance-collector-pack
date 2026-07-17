// checkimage runs the collector's own signature / SBOM / SLSA-provenance
// detection against one or more image references and prints the result as JSON.
//
// It calls the exact internal/verify code paths the collector uses at scan
// time, so it validates supply-chain detection directly against a registry -
// no cluster, no CronJob, and no dashboard UI. Use it to confirm that an image
// the collector *should* recognize as signed/attested actually comes back that
// way (for example, when a build pipeline changes how it attaches SBOM and
// provenance).
//
// Usage:
//
//	go run ./hack/checkimage <image-ref> [<image-ref> ...]
//
// Example:
//
//	go run ./hack/checkimage \
//	  quay.io/nebari/provenance-collector:0.1.0 \
//	  ghcr.io/nebari-dev/provenance-collector-pack:0.1.0
//
// If PROVENANCE_COSIGN_PUBLIC_KEY is set, signature verification uses that key;
// otherwise it only checks for signature existence (matching the collector's
// default behaviour).
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nebari-dev/provenance-collector/internal/report"
	"github.com/nebari-dev/provenance-collector/internal/verify"
)

type result struct {
	Image      string                 `json:"image"`
	Signature  *report.SignatureInfo  `json:"signature"`
	SBOM       *report.SBOMInfo       `json:"sbom"`
	Provenance *report.ProvenanceInfo `json:"provenance"`
}

func main() {
	images := os.Args[1:]
	if len(images) == 0 {
		fmt.Fprintln(os.Stderr, "usage: checkimage <image-ref> [<image-ref> ...]")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	sigVerifier := verify.NewSignatureVerifier(os.Getenv("PROVENANCE_COSIGN_PUBLIC_KEY"))
	sbomDiscoverer := verify.NewSBOMDiscoverer()
	provChecker := verify.NewProvenanceChecker()

	results := make([]result, 0, len(images))
	for _, img := range images {
		sig, err := sigVerifier.Verify(ctx, img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: signature verify %s: %v\n", img, err)
		}
		sbom, err := sbomDiscoverer.Discover(ctx, img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: sbom discover %s: %v\n", img, err)
		}
		prov, err := provChecker.Check(ctx, img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: provenance check %s: %v\n", img, err)
		}
		results = append(results, result{Image: img, Signature: sig, SBOM: sbom, Provenance: prov})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		fmt.Fprintln(os.Stderr, "encode:", err)
		os.Exit(1)
	}
}

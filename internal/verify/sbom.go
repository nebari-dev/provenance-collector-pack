package verify

import (
	"context"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	ociremote "github.com/sigstore/cosign/v2/pkg/oci/remote"

	"github.com/nebari-dev/provenance-collector/internal/report"
)

// SBOMDiscoverer checks whether container images have attached SBOM attestations.
type SBOMDiscoverer interface {
	Discover(ctx context.Context, imageRef string) (*report.SBOMInfo, error)
}

// OCISBOMDiscoverer looks for SBOM attestations in OCI registries.
type OCISBOMDiscoverer struct{}

// NewSBOMDiscoverer creates an SBOMDiscoverer that checks OCI registries.
func NewSBOMDiscoverer() SBOMDiscoverer {
	return &OCISBOMDiscoverer{}
}

func (d *OCISBOMDiscoverer) Discover(ctx context.Context, imageRef string) (*report.SBOMInfo, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return &report.SBOMInfo{}, nil
	}

	// Get the signed entity from the registry
	se, err := ociremote.SignedEntity(ref)
	if err != nil {
		return &report.SBOMInfo{}, nil
	}

	// Check for attached attestations (SBOMs are typically stored as in-toto attestations)
	atts, err := se.Attestations()
	if err != nil {
		return &report.SBOMInfo{}, nil
	}

	attList, err := atts.Get()
	if err != nil || len(attList) == 0 {
		return &report.SBOMInfo{HasSBOM: false}, nil
	}

	// Check attestation predicates for known SBOM types
	for _, att := range attList {
		payload, err := att.Payload()
		if err != nil {
			continue
		}
		format := detectSBOMFormat(string(payload))
		if format != "" {
			return &report.SBOMInfo{
				HasSBOM: true,
				Format:  format,
			}, nil
		}
	}

	return &report.SBOMInfo{HasSBOM: false}, nil
}

func detectSBOMFormat(payload string) string {
	switch {
	case strings.Contains(payload, "https://spdx.dev/Document") || strings.Contains(payload, "SPDXRef-"):
		return "spdx"
	case strings.Contains(payload, "CycloneDX") || strings.Contains(payload, "cyclonedx"):
		return "cyclonedx"
	default:
		return ""
	}
}

package verify

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	ociremote "github.com/sigstore/cosign/v2/pkg/oci/remote"

	"github.com/nebari-dev/provenance-collector/internal/report"
)

// SignatureVerifier checks whether container images have cosign signatures.
type SignatureVerifier interface {
	Verify(ctx context.Context, imageRef string) (*report.SignatureInfo, error)
}

// CosignVerifier uses the cosign library to verify image signatures.
type CosignVerifier struct {
	publicKey string
}

// NewSignatureVerifier creates a SignatureVerifier.
// If publicKey is empty, it only checks for signature existence without
// verifying the trust chain.
func NewSignatureVerifier(publicKey string) SignatureVerifier {
	return &CosignVerifier{publicKey: publicKey}
}

func (v *CosignVerifier) Verify(ctx context.Context, imageRef string) (*report.SignatureInfo, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return &report.SignatureInfo{
			Error: fmt.Sprintf("invalid image reference: %v", err),
		}, nil
	}

	// Get the signed entity from the registry
	se, err := ociremote.SignedEntity(ref)
	if err != nil {
		return &report.SignatureInfo{
			Error: fmt.Sprintf("fetching signed entity: %v", err),
		}, nil
	}

	// Check if signatures exist
	sigs, err := se.Signatures()
	if err != nil {
		return &report.SignatureInfo{
			Error: fmt.Sprintf("checking signatures: %v", err),
		}, nil
	}

	sigList, err := sigs.Get()
	if err != nil || len(sigList) == 0 {
		return &report.SignatureInfo{Signed: false}, nil
	}

	info := &report.SignatureInfo{Signed: true}

	// Full key-based verification would go here. For now, we report existence.
	if v.publicKey != "" {
		info.Verified = false
		info.Error = "key-based verification not yet implemented; signature exists"
	}

	return info, nil
}

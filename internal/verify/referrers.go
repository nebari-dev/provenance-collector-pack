package verify

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Annotation keys that carry the in-toto predicate type of a referring
// attestation manifest. Sigstore bundles use the first; buildx / BuildKit
// attestation manifests use the second.
const (
	annoSigstorePredicateType = "dev.sigstore.bundle.predicateType"
	annoInTotoPredicateType   = "in-toto.io/predicate-type"
)

// referrerManifests resolves imageRef, then reads the OCI referrers index via
// the <algo>-<hex> fallback tag derived from the image digest, and returns the
// referring manifest descriptors. Attestations attached by cosign (bundle
// format) and by docker/build-push-action (sbom/provenance) both land here.
//
// Every failure (bad ref, unreachable registry, no referrers tag) returns an
// empty slice and a nil error: a missing referrers index is the common case,
// not an error worth surfacing to the caller.
func referrerManifests(ctx context.Context, imageRef string) ([]v1.Descriptor, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, nil
	}

	desc, err := remote.Get(ref, remote.WithContext(ctx))
	if err != nil {
		return nil, nil
	}

	referrersTag := strings.Replace(desc.Digest.String(), ":", "-", 1)
	referrersRef, err := name.ParseReference(fmt.Sprintf("%s:%s", ref.Context().String(), referrersTag))
	if err != nil {
		return nil, nil
	}

	idx, err := remote.Index(referrersRef, remote.WithContext(ctx))
	if err != nil {
		return nil, nil
	}

	manifest, err := idx.IndexManifest()
	if err != nil || manifest == nil {
		return nil, nil
	}

	return manifest.Manifests, nil
}

// predicateTypes returns the candidate in-toto predicate types advertised by a
// referring manifest descriptor: the dedicated sigstore/in-toto annotations
// first, then every other annotation value as a fallback (some producers stash
// the predicate type under non-standard keys).
func predicateTypes(d v1.Descriptor) []string {
	var out []string
	if pt := d.Annotations[annoSigstorePredicateType]; pt != "" {
		out = append(out, pt)
	}
	if pt := d.Annotations[annoInTotoPredicateType]; pt != "" {
		out = append(out, pt)
	}
	for k, v := range d.Annotations {
		if k == annoSigstorePredicateType || k == annoInTotoPredicateType {
			continue
		}
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

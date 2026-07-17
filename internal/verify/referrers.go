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

// Annotation key and value that BuildKit uses to mark the attestation
// manifests it embeds in a multi-arch image index. Such entries carry platform
// unknown/unknown and reference the image manifest they attest to.
const (
	annoDockerReferenceType = "vnd.docker.reference.type"
	attestationManifestType = "attestation-manifest"
)

// referrerManifests resolves imageRef, then reads the OCI referrers index via
// the <algo>-<hex> fallback tag derived from the image digest, and returns the
// referring manifest descriptors. Cosign/sigstore bundle attestations land
// here. BuildKit's in-index attestation manifests do NOT — those are walked
// separately by indexAttestationPredicateTypes.
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

// indexAttestationPredicateTypes returns the in-toto predicate types attached
// to imageRef via BuildKit attestation manifests embedded in the image index.
//
// docker/build-push-action (sbom: true / provenance: true) and Docker Official
// Images store SBOM and SLSA attestations as extra manifests inside the
// multi-arch image index, marked with platform unknown/unknown and the
// annotation vnd.docker.reference.type=attestation-manifest. The predicate type
// lives on the attestation manifest's *layer* descriptors, so we fetch each
// attestation manifest and read its layer annotations. These attestations never
// appear in the <algo>-<hex> referrers fallback tag, so referrerManifests alone
// misses them.
//
// Every failure returns an empty slice: a single-arch image (no index) or an
// image without attestations is the common case, not an error worth surfacing.
func indexAttestationPredicateTypes(ctx context.Context, imageRef string) []string {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return nil
	}

	desc, err := remote.Get(ref, remote.WithContext(ctx))
	if err != nil || !desc.MediaType.IsIndex() {
		return nil
	}

	idx, err := desc.ImageIndex()
	if err != nil {
		return nil
	}
	manifest, err := idx.IndexManifest()
	if err != nil || manifest == nil {
		return nil
	}

	var out []string
	for _, m := range manifest.Manifests {
		if m.Annotations[annoDockerReferenceType] != attestationManifestType {
			continue
		}
		attRef := ref.Context().Digest(m.Digest.String())
		img, err := remote.Image(attRef, remote.WithContext(ctx))
		if err != nil {
			continue
		}
		attManifest, err := img.Manifest()
		if err != nil || attManifest == nil {
			continue
		}
		for _, layer := range attManifest.Layers {
			if pt := layer.Annotations[annoInTotoPredicateType]; pt != "" {
				out = append(out, pt)
			}
		}
	}
	return out
}

// artifactType prefix of a cosign/sigstore signature stored in the modern
// bundle format. Keyless `cosign sign` (v3+) attaches the signature as an OCI
// referrer with this artifactType (e.g. ...bundle.v0.3+json) rather than at the
// legacy sha256-<digest>.sig tag.
const sigstoreBundleArtifactType = "application/vnd.dev.sigstore.bundle"

// hasSignatureReferrer reports whether imageRef has a cosign/sigstore signature
// attached via the OCI referrers API. This is the modern bundle-format
// signature, which ociremote.SignedEntity (the legacy .sig tag) does not see.
//
// Uses remote.Referrers (the real referrers API) rather than the <algo>-<hex>
// fallback tag, since a registry may serve referrers only through the API.
// Any failure (bad ref, unreachable registry, no referrers) returns false: an
// image without a signature referrer is the common case, not an error.
func hasSignatureReferrer(ctx context.Context, imageRef string) bool {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return false
	}
	desc, err := remote.Get(ref, remote.WithContext(ctx))
	if err != nil {
		return false
	}
	idx, err := remote.Referrers(ref.Context().Digest(desc.Digest.String()), remote.WithContext(ctx))
	if err != nil {
		return false
	}
	manifest, err := idx.IndexManifest()
	if err != nil || manifest == nil {
		return false
	}
	for _, m := range manifest.Manifests {
		if strings.HasPrefix(string(m.ArtifactType), sigstoreBundleArtifactType) {
			return true
		}
	}
	return false
}

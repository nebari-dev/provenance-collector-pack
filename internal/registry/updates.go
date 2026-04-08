package registry

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/nebari-dev/provenance-collector/internal/report"
)

// UpdateChecker checks whether newer versions of container images are available.
type UpdateChecker interface {
	Check(ctx context.Context, imageRef string) (*report.UpdateInfo, error)
}

// RegistryUpdateChecker checks registries for newer image tags.
type RegistryUpdateChecker struct {
	opts []crane.Option
}

// NewUpdateChecker creates an UpdateChecker that queries container registries.
func NewUpdateChecker(opts ...crane.Option) UpdateChecker {
	return &RegistryUpdateChecker{opts: opts}
}

func (c *RegistryUpdateChecker) Check(ctx context.Context, imageRef string) (*report.UpdateInfo, error) {
	repo, tag := parseImageRef(imageRef)
	if tag == "" || tag == "latest" {
		return &report.UpdateInfo{CurrentTag: tag}, nil
	}

	currentVer, err := semver.NewVersion(tag)
	if err != nil {
		// Non-semver tag (e.g., sha, commit hash) — can't compare
		return &report.UpdateInfo{CurrentTag: tag}, nil
	}

	opts := append([]crane.Option{
		crane.WithContext(ctx),
		crane.WithTransport(remote.DefaultTransport),
	}, c.opts...)

	tags, err := crane.ListTags(repo, opts...)
	if err != nil {
		return nil, fmt.Errorf("listing tags for %s: %w", repo, err)
	}

	var versions []*semver.Version
	for _, t := range tags {
		v, err := semver.NewVersion(t)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}

	sort.Sort(semver.Collection(versions))

	info := &report.UpdateInfo{
		CurrentTag: tag,
	}

	// Find latest in same major version line
	for i := len(versions) - 1; i >= 0; i-- {
		v := versions[i]
		if v.Major() == currentVer.Major() && v.GreaterThan(currentVer) {
			info.LatestInMajor = v.Original()
			break
		}
	}

	// Absolute newest
	if len(versions) > 0 {
		newest := versions[len(versions)-1]
		if newest.GreaterThan(currentVer) {
			info.NewestAvailable = newest.Original()
			info.UpdateAvailable = true
		}
	}

	if info.LatestInMajor != "" {
		info.UpdateAvailable = true
	}

	return info, nil
}

// parseImageRef splits an image reference into repository and tag.
func parseImageRef(ref string) (string, string) {
	// Handle digest references (repo@sha256:...)
	if strings.Contains(ref, "@") {
		parts := strings.SplitN(ref, "@", 2)
		return parts[0], ""
	}

	// Handle tag references (repo:tag)
	lastColon := strings.LastIndex(ref, ":")
	lastSlash := strings.LastIndex(ref, "/")
	if lastColon > lastSlash && lastColon != -1 {
		return ref[:lastColon], ref[lastColon+1:]
	}

	return ref, "latest"
}

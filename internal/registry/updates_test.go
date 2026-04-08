package registry

import "testing"

func TestParseImageRef(t *testing.T) {
	tests := []struct {
		ref      string
		wantRepo string
		wantTag  string
	}{
		{"nginx:1.27", "nginx", "1.27"},
		{"nginx", "nginx", "latest"},
		{"docker.io/library/nginx:1.27-alpine", "docker.io/library/nginx", "1.27-alpine"},
		{"ghcr.io/org/repo:v2.0.0", "ghcr.io/org/repo", "v2.0.0"},
		{"myregistry.com:5000/myimage:tag", "myregistry.com:5000/myimage", "tag"},
		{"nginx@sha256:abc123", "nginx", ""},
		{"ghcr.io/org/repo@sha256:abc", "ghcr.io/org/repo", ""},
	}

	for _, tc := range tests {
		repo, tag := parseImageRef(tc.ref)
		if repo != tc.wantRepo || tag != tc.wantTag {
			t.Errorf("parseImageRef(%q) = (%q, %q), want (%q, %q)",
				tc.ref, repo, tag, tc.wantRepo, tc.wantTag)
		}
	}
}

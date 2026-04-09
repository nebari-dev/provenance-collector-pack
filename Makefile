VERSION ?= dev
IMAGE   ?= ghcr.io/nebari-dev/provenance-collector
TAG     ?= $(VERSION)

LDFLAGS := -s -w -X github.com/nebari-dev/provenance-collector/internal/report.Version=$(VERSION)

.PHONY: build test lint docker-build docker-push helm-lint clean

build:
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/provenance-collector ./cmd/provenance-collector
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/dashboard ./cmd/dashboard

test:
	go test -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...
	helm lint chart/

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(TAG) .

docker-push: docker-build
	docker push $(IMAGE):$(TAG)

helm-lint:
	helm lint chart/
	helm template test chart/ | kubeconform -strict -ignore-missing-schemas

clean:
	rm -rf bin/ coverage.out

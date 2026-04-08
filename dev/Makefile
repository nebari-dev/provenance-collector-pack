# Local development helpers using kind

CLUSTER_NAME ?= provenance-dev

.PHONY: cluster-create cluster-delete build-load deploy test-run logs report down

cluster-create:
	kind create cluster --name $(CLUSTER_NAME)

cluster-delete:
	kind delete cluster --name $(CLUSTER_NAME)

build-load:
	docker build --build-arg VERSION=dev -t provenance-collector:dev ..
	kind load docker-image provenance-collector:dev --name $(CLUSTER_NAME)

deploy: build-load
	helm upgrade --install pc ../chart/ \
		--set image.repository=provenance-collector \
		--set image.tag=dev \
		--set image.pullPolicy=Never \
		--set config.reportOutput=configmap \
		--set persistence.enabled=false \
		--set config.verifySignatures=false \
		--set config.checkUpdates=false

test-run:
	kubectl delete job manual-run --ignore-not-found
	kubectl create job --from=cronjob/pc-provenance-collector manual-run
	kubectl wait --for=condition=complete job/manual-run --timeout=120s
	kubectl logs job/manual-run

logs:
	kubectl logs job/manual-run

report:
	kubectl get configmap provenance-report -o jsonpath='{.data.report\.json}' | jq .

down:
	kind delete cluster --name $(CLUSTER_NAME)

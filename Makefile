
CONTAINER_TOOL ?= docker
KUBE_TOOL ?= kind
IMG ?= happy-beer-api:latest
NAMESPACE ?= happy-beer

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: 
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet 
	go test ./... -json -coverprofile cover.out

.PHONY: build
build: fmt vet 
	go build -o bin/api cmd/api/main.go

.PHONY: run
run:  fmt vet 
	go run ./cmd/api/main.go

.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	$(CONTAINER_TOOL) build -t ${IMG}  -f build/Dockerfile . 

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

.PHONY: 
docker-load: docker-build ## Load the docker image in the local cluster
	kind load docker-image ${IMG}

.PHONY:
deploy: docker-load ## Install the deployment in the cluster
	helm upgrade  --create-namespace -n ${NAMESPACE} happy-beer-api ./chart

# kubectl get ns ${NAMESPACE} || kubectl create ns ${NAMESPACE}
# kubectl get deployment  ${DEPLOYMENT_NAME} -n ${NAMESPACE} ||kubectl create deployment ${DEPLOYMENT_NAME} --image=${IMG} -n ${NAMESPACE} 
# kubectl patch deployments.apps  ${DEPLOYMENT_NAME} -p '{"spec": {"template": {"spec":{"containers":[{"name":"${DEPLOYMENT_NAME}","imagePullPolicy":"IfNotPresent"}]}}}}' -n ${NAMESPACE} 
#sleep 0.2 && kubectl scale --replicas=1 deployment/${DEPLOYMENT_NAME} -n ${NAMESPACE} 

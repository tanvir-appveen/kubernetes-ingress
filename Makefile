VERSION = edge
TAG = $(VERSION)
PREFIX = nginx/nginx-ingress
NGINX_VERSION = 1.19.3
GOLANG_CONTAINER = golang:1.15-alpine
GOFLAGS ?= -mod=vendor
TARGET ?= local

DOCKER_BUILD_OPTIONS = --build-arg IC_VERSION=$(VERSION)-$(GIT_COMMIT) --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg VERSION=$(VERSION) --build-arg GOLANG_CONTAINER=$(GOLANG_CONTAINER) --build-arg NGINX_VERSION=$(NGINX_VERSION)
DOCKER_CMD = docker build $(DOCKER_BUILD_OPTIONS) --target $(TARGET) -f build/Dockerfile -t $(PREFIX):$(TAG) .

GIT_COMMIT = $(shell git rev-parse --short HEAD)

export DOCKER_BUILDKIT = 1

.DEFAULT_GOAL:=help

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: test lint verify-codegen update-crds build debian-image

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: test
test: ## Run tests
	GO111MODULE=on GOFLAGS='$(GOFLAGS)' go test ./...

.PHONY: verify-codegen
verify-codegen: ## Verify code generation
	./hack/verify-codegen.sh

.PHONY: update-codegen
update-codegen: ## Generate code
	./hack/update-codegen.sh

.PHONY: update-crds
update-crds: ## Update CRDs
	go run sigs.k8s.io/controller-tools/cmd/controller-gen crd:crdVersions=v1 schemapatch:manifests=./deployments/common/crds/ paths=./pkg/apis/configuration/... output:dir=./deployments/common/crds
	go run sigs.k8s.io/controller-tools/cmd/controller-gen crd:crdVersions=v1beta1,preserveUnknownFields=false schemapatch:manifests=./deployments/common/crds-v1beta1/ paths=./pkg/apis/configuration/... output:dir=./deployments/common/crds-v1beta1
	@cp -Rp deployments/common/crds-v1beta1/ deployments/helm-chart/crds

.PHONY: certificate-and-key
certificate-and-key: ## Create default cert and key
	./build/generate_default_cert_and_key.sh

.PHONY: build
build: ## Build Ingress Controller binary
	CGO_ENABLED=0 GO111MODULE=on GOFLAGS='$(GOFLAGS)' GOOS=linux go build -installsuffix cgo -ldflags "-w -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o nginx-ingress github.com/nginxinc/kubernetes-ingress/cmd/nginx-ingress

.PHONY: debian-image
debian-image: ## Create Docker image for Ingress Controller (debian)
	$(DOCKER_CMD) --build-arg BUILD_OS=debian

.PHONY: verify-codegen
alpine-image: ## Create Docker image for Ingress Controller (alpine)
	$(DOCKER_CMD) --build-arg BUILD_OS=alpine

.PHONY: debian-image-plus
debian-image-plus: ## Create Docker image for Ingress Controller (nginx plus)
	$(DOCKER_CMD) --build-arg BUILD_OS=debian-plus --build-arg PLUS=-plus --secret id=nginx-repo.crt,src=nginx-repo.crt --secret id=nginx-repo.key,src=nginx-repo.key

.PHONY: debian-image-nap-plus
debian-image-nap-plus: ## Create Docker image for Ingress Controller (nginx plus with nap)
	$(DOCKER_CMD) --build-arg BUILD_OS=debian-plus-ap --build-arg PLUS=-plus --secret id=nginx-repo.crt,src=tempdir/nginx-repo.crt --secret id=nginx-repo.key,src=tempdir/nginx-repo.key

.PHONY: openshift-image
openshift-image: ## Create Docker image for Ingress Controller (openshift)
	$(DOCKER_CMD) --build-arg BUILD_OS=openshift

.PHONY: openshift-image-plus
openshift-image-plus: ## Create Docker image for Ingress Controller (openshift with plus)
	$(DOCKER_CMD) --build-arg BUILD_OS=openshift-plus --build-arg PLUS=-plus --secret id=nginx-repo.crt,src=tempdir/nginx-repo.crt --secret id=nginx-repo.key,src=tempdir/nginx-repo.key --secret id=rhel_license,src=tempdir/rhel_license

.PHONY: openshift-image-nap-plus
openshift-image-nap-plus: ## Create Docker image for Ingress Controller (openshift with plus and nap)
	docker build $(DOCKER_BUILD_OPTIONS) --target $(TARGET) --secret id=nginx-repo.crt,src=tempdir/nginx-repo.crt --secret id=nginx-repo.key,src=tempdir/nginx-repo.key --secret id=rhel_license,src=tempdir/rhel_license -f build/DockerfileWithAppProtectForPlusForOpenShift -t $(PREFIX):$(TAG) .

.PHONY: debian-image-opentracing
debian-image-opentracing: ## Create Docker image for Ingress Controller (with opentracing)
	docker build $(DOCKER_BUILD_OPTIONS) --target $(TARGET) -f build/DockerfileWithOpentracing -t $(PREFIX):$(TAG) .

.PHONY: debian-image-opentracing-plus
debian-image-opentracing-plus: ## Create Docker image for Ingress Controller (with opentracing and plus)
	docker build $(DOCKER_BUILD_OPTIONS) --target $(TARGET) --secret id=nginx-repo.crt,src=tempdir/nginx-repo.crt --secret id=nginx-repo.key,src=tempdir/nginx-repo.key -f build/DockerfileWithOpentracingForPlus -t $(PREFIX):$(TAG) .

.PHONY: push
push: ## Docker push to $PREFIX and $TAG
	docker push $(PREFIX):$(TAG)

.PHONY: clean
clean:  ## Remove nginx-ingress binary and temp folder
	rm -f nginx-ingress
	rm -rf tempdir

.PHONY: deps
deps: ## Add missing and remove unused modules, verify deps and make a vendored copy
	@go mod tidy && go mod verify && go mod vendor

.PHONY: clean-cache
clean-cache: ## Clean go cache
	@go clean -modcache

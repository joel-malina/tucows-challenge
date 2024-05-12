SERVICE=order-service
SERVICE_PATH=localhost/$(SERVICE)
SERVICE_BINARY_NAME=orderservice
BUILDER_IMAGE=$(SERVICE)-builder
export BUILD_DATE?=unknown
export GIT_HASH?=unknown
export GOPATH?=/tmp/gocache

RUN=docker run --rm \
	-v $(CURDIR):/opt/go/src/$(SERVICE_PATH) \
	-v $(GOPATH)/pkg/mod:/opt/go/pkg/mod \
	-w /opt/go/src/$(SERVICE_PATH) \
	-e GO111MODULE=on
.PHONY: build test clean dev lint proto

build:
	docker build -f Dockerfile.build -t $(BUILDER_IMAGE) .
	$(RUN) -e CGO_ENABLED=1 -e GOOS=linux $(BUILDER_IMAGE) \
		/bin/bash -c "git config --global --add safe.directory /opt/go/src/$(SERVICE_PATH) && \
		go build -o $(SERVICE_BINARY_NAME) -tags musl -ldflags \"-X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)\" \
		./cmd/$(SERVICE_BINARY_NAME)/..."
	$(RUN) $(BUILDER_IMAGE) abgocyclo -total -value -exclude vendor -exclude .idea . > cyclo.txt
	$(RUN) $(BUILDER_IMAGE) rhash --sha256 $(SERVICE_BINARY_NAME) -o $(SERVICE_BINARY_NAME).sha256
	docker build --tag="$(SERVICE):$(VERSION)" --tag="$(SERVICE):latest" .
	docker build --tag="$(SERVICE)-db-migration:$(VERSION)" --tag="$(SERVICE):latest" db/

clean:
	docker build -f Dockerfile.build -t $(BUILDER_IMAGE) .
	-$(RUN) $(BUILDER_IMAGE) rm -rf vendor.* *.log $(SERVICE_BINARY_NAME) $(SERVICE_BINARY_NAME).sha256 test.xml dependencies.txt
	-$(RUN) $(BUILDER_IMAGE) go clean -i
	-$(RUN) $(BUILDER_IMAGE) find . -type f -name 'coverage.xml' -delete
	-docker rmi -f $(SERVICE):$(VERSION) $(SERVICE):latest $(BUILDER_IMAGE)

test:
	echo "not running tests"

dev:
	docker-compose up

cleandev:
	docker-compose down
	-docker volume remove tucows-challenge_postgres-vol
	docker-compose up -d
	sleep 2
	db/apply_local.sh

lint:
	make lint-golang-ci

lint-golang-ci:
	docker run --rm -v $(CURDIR):$(CURDIR) -w $(CURDIR) golangci/golangci-lint:v1.54.2 golangci-lint run -v --fix
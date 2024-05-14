ORDER_SERVICE=order-service
ORDER_SERVICE_PATH=localhost/$(ORDER_SERVICE)
ORDER_BINARY_NAME=order
ORDER_BUILDER_IMAGE=$(ORDER_SERVICE)-builder
PAYMENT_SERVICE=payment-service
PAYMENT_SERVICE_PATH=localhost/$(PAYMENT_SERVICE)
PAYMENT_BINARY_NAME=payment
PAYMENT_BUILDER_IMAGE=$(PAYMENT_SERVICE)-builder

export BUILD_DATE?=unknown
export GIT_HASH?=unknown
export GOPATH?=/tmp/gocache
export VERSION=demo
CURDIR:=$(PWD)

RUN=docker run --rm \
	-v $(CURDIR):/opt/go/src/$(SERVICE_PATH) \
	-v $(GOPATH)/pkg/mod:/opt/go/pkg/mod \
	-w /opt/go/src/$(SERVICE_PATH) \
	-e GO111MODULE=on
.PHONY: build test clean dev lint proto build-order build-payment

build-order:
	@echo "Building Order Service..."
	docker build -f Dockerfile.build -t $(ORDER_BUILDER_IMAGE) .
	$(RUN) -e CGO_ENABLED=1 -e GOOS=linux $(ORDER_BUILDER_IMAGE) \
		/bin/sh -c "git config --global --add safe.directory /opt/go/src/$(ORDER_SERVICE_PATH) && \
		go build -o $(ORDER_BINARY_NAME) -tags musl -ldflags \"-X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)\" \
		-buildvcs=false ./cmd/$(ORDER_BINARY_NAME)/..."
	$(RUN) $(ORDER_BUILDER_IMAGE) rhash --sha256 $(ORDER_BINARY_NAME) -o $(ORDER_BINARY_NAME).sha256
	docker build --tag="$(ORDER_SERVICE):$(VERSION)" --tag="$(ORDER_SERVICE):latest" .
	docker build --tag="$(ORDER_SERVICE)-db-migration:$(VERSION)" --tag="$(ORDER_SERVICE):latest" db/

build-payment:
	@echo "Building Payment Service..."
	docker build -f Dockerfile.build -t $(PAYMENT_BUILDER_IMAGE) .
	$(RUN) -e CGO_ENABLED=1 -e GOOS=linux $(PAYMENT_BUILDER_IMAGE) \
		/bin/sh -c "git config --global --add safe.directory /opt/go/src/$(PAYMENT_SERVICE_PATH) && \
		go build -o $(PAYMENT_BINARY_NAME) -tags musl -ldflags \"-X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)\" \
		-buildvcs=false ./cmd/$(PAYMENT_BINARY_NAME)/..."
	$(RUN) $(PAYMENT_BUILDER_IMAGE) rhash --sha256 $(PAYMENT_BINARY_NAME) -o $(PAYMENT_BINARY_NAME).sha256
	docker build --tag="$(PAYMENT_SERVICE):$(VERSION)" --tag="$(PAYMENT_SERVICE):latest" .

build: build-order build-payment

clean-payment:
	docker build -f Dockerfile.build -t $(PAYMENT_BUILDER_IMAGE) .
	-$(RUN) $(PAYMENT_BUILDER_IMAGE) rm -rf vendor.* *.log $(PAYMENT_BINARY_NAME) $(PAYMENT_BINARY_NAME).sha256 test.xml dependencies.txt
	-$(RUN) $(PAYMENT_BUILDER_IMAGE) go clean -i
	-docker rmi -f $(PAYMENT_SERVICE):$(VERSION) $(PAYMENT_SERVICE):latest $(PAYMENT_BUILDER_IMAGE)

clean-order:
	docker build -f Dockerfile.build -t $(ORDER_BUILDER_IMAGE) .
	-$(RUN) $(ORDER_BUILDER_IMAGE) rm -rf vendor.* *.log $(ORDER_BINARY_NAME) $(ORDER_BINARY_NAME).sha256 test.xml dependencies.txt
	-$(RUN) $(ORDER_BUILDER_IMAGE) go clean -i
	-docker rmi -f $(ORDER_SERVICE):$(VERSION) $(ORDER_SERVICE):latest $(ORDER_BUILDER_IMAGE)

clean: clean-order clean-payment

test:
	echo "not running tests"

dev:
	docker-compose up

cleandev:
	docker-compose down
	-docker volume remove tucows-challenge_postgres-vol
	-docker volume remove tucows-challenge_rabbitmq-data
	docker-compose up -d
	sleep 2
	db/apply_local.sh

lint:
	make lint-golang-ci

lint-golang-ci:
	docker run --rm -v $(CURDIR):$(CURDIR) -w $(CURDIR) golangci/golangci-lint:v1.54.2 golangci-lint run -v --fix
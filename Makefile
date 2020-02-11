NAME = chestercheng/taiko-web
VERSION = latest

GOFMT ?= gofmt "-s"
GO ?= go

PACKAGES ?= $(shell $(GO) list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)

BUILD_OPT =
BUILD_FLAGS = -X taiko-web/config.version=$(VERSION)
BUILD_FLAGS += -X taiko-web/config.commit=$(shell git rev-list -1 HEAD)
BUILD_FLAGS += -X taiko-web/config.url=$(shell git remote get-url origin | sed 's/.git$$/\//g')

fmt:
	$(GOFMT) -w $(SOURCES)

vet:
	$(GO) vet $(PACKAGES)

devserver:
	$(GO) run main.go --mode debug

build:
	docker build -t $(NAME):$(VERSION) --rm $(BUILD_OPT) \
		--build-arg BUILD_FLAGS="$(BUILD_FLAGS)" \
		.

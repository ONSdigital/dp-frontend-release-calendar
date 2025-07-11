BINPATH ?= build


GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)
LOCAL_DP_RENDERER_IN_USE = $(shell grep -c "\"github.com/ONSdigital/dp-renderer/v2\" =" go.mod)

SERVICE_PATH = github.com/ONSdigital/dp-frontend-release-calendar/service

LDFLAGS = -ldflags "-X $(SERVICE_PATH).BuildTime=$(BUILD_TIME) -X $(SERVICE_PATH).GitCommit=$(GIT_COMMIT) -X $(SERVICE_PATH).Version=$(VERSION)"

.PHONY: all
all: delimiter-AUDIT audit delimiter-LINTERS lint delimiter-UNIT-TESTS test delimiter-COMPONENT-TESTS test-component delimiter-FINISH ## Runs multiple targets, audit, lint, test and test-component

.PHONY: audit
audit: ## Runs checks for security vulnerabilities on dependencies (including transient ones)
	go list -json -m all | nancy sleuth

.PHONY: build
build: generate-prod ## Builds binary of application code and stores in bin directory as dp-frontend-release-calendar
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-frontend-release-calendar

.PHONY: convey
convey: ## Runs unit test suite and outputs results on http://127.0.0.1:8080/
	goconvey ./...

.PHONY: delimiter-%
delimiter-%:
	@echo '===================${GREEN} $* ${RESET}==================='

.PHONY: debug
debug: generate-debug ## Used to build and run code locally in debug mode
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-frontend-release-calendar
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-release-calendar

.PHONY: fmt
fmt: ## Run Go formatting on code
	go fmt ./...

.PHONY: lint
lint: generate-prod
	golangci-lint run ./... --build-tags 'production'

.PHONY: test
test: generate-prod ## Runs unit tests including checks for race conditions and returns coverage
	go test -count=1 -race -cover -tags 'production' ./...

.PHONY: test-component
test-component: generate-prod ## Runs component test suite
	go test -cover -tags 'production' -coverpkg=github.com/ONSdigital/dp-frontend-release-calendar/... -component

.PHONY: generate-debug
generate-debug: fetch-renderer-lib
	cd assets; go run github.com/kevinburke/go-bindata/go-bindata -prefix $(CORE_ASSETS_PATH)/assets -debug -o data.go -pkg assets locales/... templates/... $(CORE_ASSETS_PATH)/assets/locales/... $(CORE_ASSETS_PATH)/assets/templates/...
	{ echo "// +build debug\n"; cat assets/data.go; } > assets/debug.go.new
	mv assets/debug.go.new assets/data.go

.PHONY: generate-prod
generate-prod: fetch-renderer-lib
	cd assets; go run github.com/kevinburke/go-bindata/go-bindata -prefix $(CORE_ASSETS_PATH)/assets -o data.go -pkg assets locales/... templates/... $(CORE_ASSETS_PATH)/assets/locales/... $(CORE_ASSETS_PATH)/assets/templates/...
	{ echo "// +build production\n"; cat assets/data.go; } > assets/data.go.new
	mv assets/data.go.new assets/data.go

.PHONY: fetch-dp-renderer
fetch-renderer-lib:
ifeq ($(LOCAL_DP_RENDERER_IN_USE), 1)
	$(eval CORE_ASSETS_PATH = $(shell grep -w "\"github.com/ONSdigital/dp-renderer/v2\" =>" go.mod | awk -F '=> ' '{print $$2}' | tr -d '"'))
else
	$(eval APP_RENDERER_VERSION=$(shell grep "github.com/ONSdigital/dp-renderer/v2" go.mod | cut -d ' ' -f2 ))
	$(eval CORE_ASSETS_PATH = $(shell go get github.com/ONSdigital/dp-renderer/v2@$(APP_RENDERER_VERSION) && go list -f '{{.Dir}}' -m github.com/ONSdigital/dp-renderer/v2))
endif

.PHONY: help
help: ## Show help page for list of make targets
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

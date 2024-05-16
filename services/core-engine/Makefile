PHONY: help run test coverage check build mock

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Running code at local for testing
	@go run main.go

test: ## Running Unit-test Code
	@go test -race $$(go list ./... | grep -v /vendor/) -coverprofile coverage.out

coverage: ## Running Code Coverage
	@go fmt ./... && go test -coverprofile coverage.cov -cover ./... # use -v for verbose

check: ## Running Code Dependency Check
	@go mod tidy
	@go mod download
	@go mod verify

build: ## Build Code to Binary Artifact
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o $(CI_PROJECT_DIR)/$(ARTIFACT_DIR)/$(CI_PROJECT_NAME)

mock: ## Automatically regenerate all mocking interface
	@go get github.com/maxbrunsfeld/counterfeiter/v6
	@go generate ./...
PKG := "github.com/rocketblend/rocketblend-launcher"
PKG_LIST := $(shell go list ${PKG}/...)

.DEFAULT_GOAL := check
check: dep fmt vet test ## Check project
	
vet: ## Vet the files
	@go vet ${PKG_LIST}

fmt: ## Style check the files
	@gofmt -s -w .

test: ## Run tests
	@go test -short ${PKG_LIST}

race: ## Run tests with data race detector
	@go test -race ${PKG_LIST}

benchmark: ## Run benchmarks
	@go test -run="-" -bench=".*" ${PKG_LIST}

dep:
	@go mod download
	@go mod tidy

run:
	@go run ./cmd/launcher

build:
	@go build ./cmd/launcher

dry:
	@goreleaser release --snapshot --rm-dist

release:
	@git tag $(version)
	@git push origin $(version)
	@goreleaser --rm-dist

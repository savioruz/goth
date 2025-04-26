export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)
DB_PATH:=$(CURDIR)/database/postgres

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_.-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

deps.bin: ### install tools (mandatory for development)
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	GOBIN=$(LOCAL_BIN) go install github.com/air-verse/air@latest
	GOBIN=$(LOCAL_BIN) go install golang.org/x/vuln/cmd/govulncheck@latest
	GOBIN=$(LOCAL_BIN) go install github.com/google/wire/cmd/wire@latest
.PHONY: deps.bin

deps.audit: ### check dependencies vulnerabilities
	$(LOCAL_BIN)/govulncheck ./...
.PHONY: deps.audit

generate.domains: ### domains=$DOMAIN (generate domains including sqlc.yaml)
		@if [ -z "$(domains)" ]; then \
			echo "Please set the domains variable"; \
			exit 1; \
		fi
		mkdir -p ./internal/domains/$(domains)/service \
			./internal/domains/$(domains)/handler \
			./internal/domains/$(domains)/repository \
			$(DB_PATH)/domains/$(domains)
		touch $(DB_PATH)/domains/$(domains)/schema.sql $(DB_PATH)/domains/$(domains)/queries.sql
		@echo "version: \"2\"" > $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "sql:" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "  - name: \"$(domains)\"" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "    engine: \"postgresql\"" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "    schema: \"./schema.sql\"" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "    queries: \"./queries.sql\"" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "    gen:" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "      go:" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "        package: \"sqlc\"" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "        sql_package: \"pgx/v5\"" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "        out: \"./sqlc\"" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "        emit_json_tags: true" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "        emit_db_tags: true" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "        emit_methods_with_db_argument: true" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "        emit_interface: true" >> $(DB_PATH)/domains/$(domains)/sqlc.yaml
		@echo "Domain structure created at ./internal/domains/$(domains) and sqlc.yaml at $(DB_PATH)/domains/$(domains)"
.PHONY: generate.domains

generate.sqlc: ### domains=$DOMAIN (generate sqlc code)
	@if [ -z "$(domains)" ]; then \
		echo "Please set the domains variable"; \
		exit 1; \
	fi
	$(LOCAL_BIN)/sqlc generate --file $(DB_PATH)/domains/$(domains)/sqlc.yaml
.PHONY: generate.sqlc

generate.swag: ### generate swagger docs
	$(LOCAL_BIN)/swag init --parseDependency --parseInternal --parseDepth=2 -g ./internal/delivery/http/router.go
.PHONY: generate.swag

generate.mock: ### generate mock
	@for domain in $$(find ./internal/domains -mindepth 1 -maxdepth 1 -type d -exec basename {} \;); do \
		mkdir -p ./internal/domains/$$domain/mock; \
		for dir in repository service; do \
			if [ -d "./internal/domains/$$domain/$$dir" ]; then \
				f=$$(find "./internal/domains/$$domain/$$dir" -name "*.go" -not -path "*/mock/*" -type f | xargs grep -l "type.*interface\|type.*Interface" 2>/dev/null || true); \
				if [ -n "$$f" ]; then \
					echo "$$f" | while read file; do \
						if [ -n "$$file" ]; then \
							dest_file="./internal/domains/$$domain/mock/$$(basename $${file%.*})_mock.go"; \
							$(LOCAL_BIN)/mockgen -source="$$file" -destination="$$dest_file" -package=mock || echo "    ERROR: Failed to generate mock for $$file"; \
						fi \
					done; \
				fi; \
			fi; \
		done; \
	done
	go generate ./pkg/...
	@echo "Mock generation completed"
.PHONY: generate.mock

generate: generate.swag generate.mock ### generate code
	cd ./internal/app && go generate ./... && wire ./wire.go
.PHONY: generate

lint: ### check by golangci linter
	$(LOCAL_BIN)/golangci-lint run
.PHONY: linter-golangci

test: generate ### run test
	@if ! -d ./tmp ]; then \
		mkdir -p ./tmp; \
	fi
	go test -v -race -covermode atomic -coverprofile=tmp/coverage.txt ./internal/...
.PHONY: test

coverage: ### show coverage
	go tool cover -html=tmp/coverage.txt

dev: generate ### Run dev
	$(LOCAL_BIN)/air -c ./.air.toml
.PHONY: dev

run: deps swag-v1 ### swag run for API v1
	go mod download && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

clean: ### clean
	rm -rf ./bin ./tmp ./test/mock ./test/coverage.txt ./test/coverage.html
.PHONY: clean

.PHONY: build clean dev frontend-dist generate lint test check frontend-build verify-codegen

# Build frontend and copy to embed location
frontend-dist:
	cd frontend && pnpm build
	rm -rf api/embedded/dist
	cp -r frontend/build api/embedded/dist

# Run go generate (includes ogen code generation)
generate:
	go install github.com/ogen-go/ogen/cmd/ogen@latest
	cd api && go generate ./...

# Build the single binary with embedded frontend
build: frontend-dist generate build-api

build-api:
	cd api && go build -o ../meet-mesh ./cmd

# Clean build artifacts
clean:
	rm -rf frontend/build api/embedded/dist meet-mesh

# Run development servers
dev:
	@echo "Run 'cd frontend && pnpm dev' and 'cd api && go run ./cmd' in separate terminals"

# Lint Go code with golangci-lint
lint:
	cd api && golangci-lint run ./...

# Run Go tests with race detector
test:
	cd api && CGO_ENABLED=1 go test -race ./...

# Run svelte-check on frontend
check:
	cd frontend && pnpm check

# Build frontend only (no copy to embed dir)
frontend-build:
	cd frontend && pnpm build

# Verify generated code is up to date
verify-codegen:
	$(MAKE) generate
	cd frontend && pnpm generate:api
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: Generated code is out of date. Please run 'make generate' and 'cd frontend && pnpm generate:api' and commit the changes."; \
		git status --porcelain; \
		git diff; \
		exit 1; \
	fi

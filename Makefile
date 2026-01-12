.PHONY: build clean dev frontend-dist generate

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

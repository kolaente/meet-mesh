# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder

WORKDIR /build

ENV PNPM_CACHE_FOLDER=.cache/pnpm/

COPY frontend/pnpm-lock.yaml frontend/package.json ./
RUN npm install -g corepack && corepack enable && \
    pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

# Stage 2: Build Go backend with xgo for cross-compilation with CGO
FROM --platform=$BUILDPLATFORM ghcr.io/techknowlogick/xgo:go-1.25.x AS backend-builder

WORKDIR /go/src/github.com/kolaente/meet-mesh

# Copy go mod files first for caching
COPY api/go.mod api/go.sum ./api/
RUN cd api && go mod download

# Copy everything
COPY . ./

# Copy frontend build from previous stage
COPY --from=frontend-builder /build/build ./api/embedded/dist/

ARG TARGETOS TARGETARCH TARGETVARIANT

# Build with xgo for the target platform
RUN cd api && \
    /usr/local/bin/xgo \
        -targets "${TARGETOS}/${TARGETARCH}/${TARGETVARIANT}" \
        -dest /build \
        -out meet-mesh \
        -ldflags="-s -w" \
        ./cmd && \
    mv /build/meet-mesh-* /build/meet-mesh

# Stage 3: Final minimal image
FROM scratch

LABEL org.opencontainers.image.authors="konrad@kolaente.de"
LABEL org.opencontainers.image.source="https://github.com/kolaente/meet-mesh"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="Meet Mesh"

WORKDIR /app/meet-mesh
ENTRYPOINT ["/app/meet-mesh/meet-mesh"]
EXPOSE 8080
USER 1000

ENV MEET_MESH_DATABASE_PATH=/data/meet-mesh.db

# Copy the binary
COPY --from=backend-builder /build/meet-mesh meet-mesh

# Copy SSL certificates for HTTPS
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

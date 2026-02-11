# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder

WORKDIR /build

ENV PNPM_CACHE_FOLDER=.cache/pnpm/

COPY frontend/pnpm-lock.yaml frontend/package.json ./
RUN npm install -g corepack && corepack enable && \
    pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build


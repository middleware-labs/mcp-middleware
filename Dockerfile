# syntax=docker/dockerfile:1

# ---- build stage ---------------------------------------------------------
# Cross-compiles a static binary for the target platform (set by buildx).
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
WORKDIR /src

# Download dependencies first so this layer is cached across source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
# CGO disabled → fully static binary that runs on distroless/static.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/mcp-middleware .

# ---- runtime stage -------------------------------------------------------
# distroless/static:nonroot — no shell, no package manager, runs as uid 65532.
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /out/mcp-middleware /app/mcp-middleware

# HTTP mode, bound to all interfaces inside the container. Override at runtime
# via env (e.g. Kubernetes ConfigMap). MW_AUTH_SERVER_URL is required in http mode.
ENV APP_MODE=http \
    APP_HOST=0.0.0.0 \
    APP_PORT=8080

EXPOSE 8080
USER nonroot:nonroot
# No ENTRYPOINT/CMD: the run command is provided by the Kubernetes Deployment
# (`command: ["/app/mcp-middleware"]`). To run the image standalone, pass the binary:
#   docker run <image> /app/mcp-middleware

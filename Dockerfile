FROM golang:1.22-bookworm AS base
WORKDIR /build
ENV CGO_ENABLED=0
# Install dependencies
COPY go.* .
RUN go mod download

# Build the binary
# '--mount=target=.': use bind mounting from the build context
# '--mount=type=cache,target=/root/.cache/go-build': use go’s compiler cache
FROM base AS build
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -trimpath -ldflags "-s -w -extldflags '-static'" \
    -o /app ./cmd/app/main.go

# Run the integration tests (a script that uses go commands)
FROM base AS integration-test
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    ./scripts/integration-test.sh

FROM base AS lint
# Run the linter
# '--mount=from=golangci/golangci-lint,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint': use binary file from another image
# '--mount=type=cache,target=/root/.cache/golangci-lint': use golangci-lint cache
RUN --mount=target=. \
    --mount=from=golangci/golangci-lint,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    golangci-lint run

FROM debian:bookworm-slim AS app
# Install ffmpeg
RUN apt -y update && apt -y upgrade && apt install -y ffmpeg
# Copy the binary
COPY --from=build /app /app
# Copy Swagger documentation
COPY ./swagger/swagger.* ./swagger
# Create environment
COPY env/app.env /
# Run the binary
ENTRYPOINT ["/app", "--config=/app.env"]

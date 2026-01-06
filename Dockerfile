# Build stage
FROM golang:1.25-alpine AS builder
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ARG VERSION=dev
RUN go build -ldflags "-X main.version=${VERSION}" -o mcp-parseable-server ./cmd/mcp_parseable_server

# Final image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/mcp-parseable-server .
EXPOSE 9034
ENTRYPOINT ["./mcp-parseable-server"]
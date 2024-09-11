# Stage 1: Build stage
FROM golang:1.23.1-alpine AS builder

# Set the working directory for the build
WORKDIR /build

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the source code
COPY ./api ./api

# Copy environment file
COPY ./.env ./.env

# Set environment variables
ENV GO111MODULE=on
ENV GOCACHE=/root/.cache/go-build

# Build the Go application with CGO enabled for Alpine
RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_CFLAGS_ALLOW=-Xpreprocessor \
    GOOS=linux go build -a -installsuffix cgo -o apiserver ./api/cmd/servid

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory in the final stage
WORKDIR /app

# Copy the binary and environment file from the build stage
COPY --from=builder /build/apiserver /app/
COPY --from=builder /build/.env /app/


# Use a non-root usergrp for security
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

# Change ownership of the app directory and its contents to the nonroot usergrp
RUN chown -R nonroot:nonroot /app

USER nonroot

# Set the entrypoint to the Go application
ENTRYPOINT ["/app/apiserver"]
EXPOSE 3000

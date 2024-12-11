# Stage 1: Build stage
FROM golang:1.23.1-alpine AS builder

# Set the working directory for the build
WORKDIR /build

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the source code
COPY ./api ./api
COPY ./internal ./internal
COPY ./business ./business
COPY ./main.go ./main.go
# Copy environment file
COPY ./.env ./.env

# Set environment variables
ENV GO111MODULE=on
ENV GOCACHE=/root/.cache/go-build

# Build the Go application with CGO enabled for Alpine
RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_CFLAGS_ALLOW=-Xpreprocessor \
    GOOS=linux go build -a -installsuffix cgo -o innovia-backend .

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory in the final stage
WORKDIR /app

# Copy the binary and environment file from the build stage
COPY --from=builder /build/innovia-backend /app/
COPY --from=builder /build/.env /app/

# Use nonroot user
#RUN addgroup -S nonroot \
#    && adduser -S nonroot -G nonroot
#
#USER nonroot

# Set the entrypoint to the Go application
ENTRYPOINT ["/app/innovia-backend"]
EXPOSE 3000

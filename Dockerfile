# Use the latest version of Go as the base image
FROM golang:1.22 AS base

# Install needed dependencies for base image and update certs
RUN apt-get update \
  && apt-get install -y ca-certificates \
    && update-ca-certificates

# create a build artifact
FROM base AS builder
# Set the working directory to the root of the project
WORKDIR /app

# Copy the Go dependencies file and download the dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the Makefile and the rest of the source code
COPY .git ./.git
COPY . ./

# Build the application with static linking for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telemetry .

# Create a new, smaller image based on alpine
FROM alpine:latest

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -g 1001 telemetry && \
    adduser -D -s /bin/sh -u 1001 -G telemetry telemetry

# Make dir for mounting config file and set ownership
RUN mkdir -p /home/telemetry/.telemetry/config && \
    chown -R telemetry:telemetry /home/telemetry

# Copy the built executable from the builder image and set ownership
COPY --from=builder /app/telemetry /app/telemetry
RUN chown telemetry:telemetry /app/telemetry

# Switch to the non-root user
USER telemetry

ENTRYPOINT ["/app/telemetry"]
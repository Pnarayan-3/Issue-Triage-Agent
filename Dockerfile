# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o issue-triage-agent ./cmd/agent

# Runtime stage
FROM alpine:3.22

WORKDIR /app

# Add CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/issue-triage-agent .

ENTRYPOINT ["./issue-triage-agent"]
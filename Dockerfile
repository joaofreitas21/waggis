# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go mod files 
COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN templ generate

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and required files
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/data ./data

EXPOSE 8080

CMD ["./main"]
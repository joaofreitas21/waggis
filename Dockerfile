# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app


# Copy go mod files 
COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@$(go list -m -f '{{.Version}}' github.com/a-h/templ)


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
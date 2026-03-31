
FROM golang:1.25-alpine AS builder

WORKDIR /app


RUN apk add --no-cache nodejs npm


RUN go install github.com/a-h/templ/cmd/templ@latest


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN templ generate


COPY package.json package-lock.json ./
RUN npm ci


RUN npx tailwindcss -i ./static/input.css -o ./static/styles.css --minify


RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go



FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/data ./data

EXPOSE 8080
CMD ["./main"]
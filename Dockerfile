
# Tailwind CSS Build
FROM node:slim AS assets
WORKDIR /app
COPY package.json .
COPY package-lock.json .
RUN npm ci

COPY input.css .
RUN mkdir -p public \
    && npx @tailwindcss/cli -i ./input.css -o ./public/tailwind.css --minify

# Web Server Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN templ generate
COPY --from=assets /app/public ./public
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api.go

# Final Image
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/public ./public
EXPOSE 1323
ENTRYPOINT ["./server"]

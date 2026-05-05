# docker/Dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy the go mod files and download dependencies
COPY go.mod go.sum ./
RUN touch go.sum && go mod download

# Copy the rest of the monorepo
COPY . .

# Build the specified bot
ARG BOT_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/${BOT_NAME}

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bot .

# Execute the bot
CMD ["./bot"]

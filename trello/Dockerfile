FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o trello-clone-backend main.go

# --- Final Stage ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy built binary from builder stage
COPY --from=builder /app/trello-clone-backend .
# Copy .env.example, actual .env should be mounted or vars provided
COPY --from=builder /app/.env.example .

# Standard port
EXPOSE 8080

# Run the executable
CMD ["./trello-clone-backend"]
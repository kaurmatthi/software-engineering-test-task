# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cruder ./cmd

# Minimal runtime image
FROM alpine:latest

RUN addgroup -S crudergroup && adduser -S cruder -G crudergroup

WORKDIR /root/
COPY --from=builder /app/cruder .

USER cruder

EXPOSE 8080

CMD ["./cruder"]
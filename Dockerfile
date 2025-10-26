#for builder
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a \
    -ldflags='-w -s -extldflags "-static"' \
    -o main .
#for running
FROM alpine:3.18
RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder --chown=appuser:appgroup /app/main .

USER appuser
EXPOSE 8080
CMD ["./main"]

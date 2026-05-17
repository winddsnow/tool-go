FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tool-go .

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/tool-go .
COPY --from=builder /app/manifest/config ./manifest/config
COPY --from=builder /app/resource ./resource

EXPOSE 8000

CMD ["./tool-go"]

FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN go build -o sidecar .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sidecar .
CMD ["./sidecar"]

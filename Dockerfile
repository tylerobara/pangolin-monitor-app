FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o pangolin-monitor main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/pangolin-monitor .
ENTRYPOINT ["./pangolin-monitor"]
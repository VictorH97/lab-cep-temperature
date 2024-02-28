FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C cmd/server -o ../../server

FROM scratch
WORKDIR /app
COPY --from=builder /app/server .
ENTRYPOINT ["./server"]
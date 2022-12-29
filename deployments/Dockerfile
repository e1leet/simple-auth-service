FROM golang:1.19-alpine
WORKDIR /wishit
COPY ../.. .
RUN go build -o app ./cmd/app/main.go
CMD ["./app", "-config=configs/config.local.yaml"]

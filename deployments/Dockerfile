FROM golang:1.19-alpine
WORKDIR /auth-service
COPY ./.. .
RUN go build -o app ./cmd/app/main.go
CMD ["./app", "-config=configs/config.local.yaml"]

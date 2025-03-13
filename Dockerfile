FROM golang:1.23.3-alpine as builder

COPY . /github.com/Muvi7z/chat-auth-s/
WORKDIR /github.com/Muvi7z/chat-auth-s/

RUN go mod download
RUN gp build -o ./bin cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Muvi7z/chat-auth-s/bin .

CMD ["./bin"]
FROM golang:1.24-alpine

WORKDIR /app
COPY . .
RUN go mod init consumer && go get github.com/segmentio/kafka-go && go build -o consumer
CMD ["./consumer"]

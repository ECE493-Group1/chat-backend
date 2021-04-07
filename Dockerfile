FROM golang:1.14-alpine

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./cmd

CMD ["./main"]
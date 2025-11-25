FROM golang:1.23-alpine

RUN apk add --no-cache git openssh

WORKDIR /app

COPY . .

RUN go mod tidy


RUN go build -o go-article ./cmd/app

EXPOSE 8002

CMD ["./go-article"]

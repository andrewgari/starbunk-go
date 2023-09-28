FROM golang:alpine

RUN apk add --no-cache git

WORKDIR /etc/github.com/andrewgari/starbunk-go

RUN go mod init github.com/andrewgari/starbunk-go

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

CMD ["go", "run", "."]
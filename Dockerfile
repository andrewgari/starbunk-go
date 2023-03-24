FROM golang:latest

RUN apt update
WORKDIR /etc/github.com/andrewgari/starbunk-go
COPY go.mod .
COPY go.sum .

RUN go mod download
ADD . .
RUN go get -d -v ./...
RUN go install -v ./...
# RUN go build -v -o /usr/local/bin/src ./..
CMD ["/go/bin/starbunk-go"]
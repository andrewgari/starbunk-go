FROM golang:latest
WORKDIR $GOPATH/src/github.com/covadax1/starbunk-go
COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
# RUN go build -v -o /usr/local/bin/src ./..

CMD ["starbunk-go"]
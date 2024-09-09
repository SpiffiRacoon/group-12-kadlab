FROM golang:1.23-alpine

WORKDIR /kademlia-app

COPY go.mod .
COPY main.go ./
COPY kademlia/*.go ./kademlia/

RUN go mod download

RUN go build -o main .

CMD ["./main"]
FROM golang:1.23-alpine

WORKDIR /kademlia-app

COPY go.mod .
COPY main.go ./
COPY kademlia/*.go ./kademlia/
COPY kademlia/cli/*.go ./kademlia/cli/

RUN go mod download

RUN go build -o main .

# We need curl becaus alpine comes with NOTHING!!!!!
RUN apk --no-cache add curl

EXPOSE 3000

CMD ["./main"]
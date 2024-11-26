FROM golang:1.23.3-alpine3.20
WORKDIR /usr/src/purgatory
#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/local/bin/purgatory ./...
CMD ["purgatory"]

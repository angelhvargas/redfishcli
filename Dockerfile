FROM golang:1.21-bookworm

WORKDIR /usr/src/redfishcli

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin ./...

CMD ["redfishcli", "-h"]
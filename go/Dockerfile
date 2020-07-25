FROM golang:1.14

WORKDIR /fiber

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./src /fiber

RUN go build -ldflags="-s -w" -o app .

CMD ["./app"]

FROM golang:1.14 as build

WORKDIR /fiber

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./src /fiber

RUN go build -ldflags="-s -w" -o app .

FROM debian:buster-slim
COPY --from=build /fiber/app /app 

CMD ["/app"]

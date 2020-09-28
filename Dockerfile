FROM golang:1-buster

CMD mkdir /app

WORKDIR /app
COPY . .
RUN go mod tidy && go build

CMD ["./plebiscitobot"]


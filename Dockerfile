FROM golang:1.20

WORKDIR /go/src/althea-api

COPY . .

RUN go build -o build/althea-api main.go

CMD ["./build/althea-api"]
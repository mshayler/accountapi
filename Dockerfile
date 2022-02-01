#Go Build Image
FROM golang:1.17

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build cmd/main.go

EXPOSE 8081

CMD ["./main"]
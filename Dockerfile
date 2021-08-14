FROM golang:1.13.10-buster
WORKDIR /go/src
COPY . .

RUN go get

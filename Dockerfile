FROM golang:1.15.2-alpine as build

WORKDIR $GOPATH/src/github.com/shekhirin/zenly-task

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /go/bin/zenly ./cmd/zenly/

FROM alpine:3.12

WORKDIR /app

COPY --from=build /go/bin/zenly ./

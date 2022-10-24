FROM golang:1.16-alpine

RUN apk update
RUN apk upgrade
RUN apk add bash
RUN apk add curl

WORKDIR /app
#ADD go.sum .

COPY go.mod .
#COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/dist .

CMD ./out/dist
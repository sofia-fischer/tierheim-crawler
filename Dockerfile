FROM golang:1.20.4

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy

ENTRYPOINT ["air"]

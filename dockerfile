FROM golang:1.21.1-alpine3.18

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o bin .
EXPOSE 8080

ENTRYPOINT [ "/app/bin" ]
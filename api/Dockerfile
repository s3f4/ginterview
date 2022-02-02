FROM golang:1.17.3-alpine3.14
RUN go get github.com/githubnemo/CompileDaemon
WORKDIR /app/
COPY go.mod .
RUN go mod download
RUN go mod verify
COPY ./api ./api
RUN go mod tidy
WORKDIR /app/api
CMD CompileDaemon -directory=. -build="go build -o api" -command="./api"


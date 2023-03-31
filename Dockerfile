FROM golang:1.20.2

WORKDIR /app
COPY ./src/go.mod ./
COPY ./src/go.sum ./
RUN go mod download
COPY ./src ./
RUN go build -o go-oauth2-server.out
CMD ./go-oauth2-server.out ${ARGS}

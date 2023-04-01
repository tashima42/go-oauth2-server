FROM golang:1.19

WORKDIR /app
COPY ./api/go.mod ./
COPY ./api/go.sum ./
RUN go mod download
COPY ./api ./
RUN go build -o go-oauth2-server.out
CMD ./go-oauth2-server.out ${ARGS}

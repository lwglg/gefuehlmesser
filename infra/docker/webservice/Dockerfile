FROM golang:1.25.1-alpine

ENV WS_ROOT=./webservice
ENV WORKDIR=/home/gefuehlmesser/webservice

WORKDIR ${WORKDIR}

RUN apk add --no-cache gcc musl-dev

COPY ${WS_ROOT}/go.mod ./go.mod
COPY ${WS_ROOT}/go.sum ./go.sum

RUN go mod download

COPY ${WS_ROOT}/api ./api
COPY ${WS_ROOT}/bin ./bin
COPY ${WS_ROOT}/cmd ./cmd
COPY ${WS_ROOT}/config ./config
COPY ${WS_ROOT}/docs ./docs
COPY ${WS_ROOT}/libs ./libs

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api

EXPOSE 8080

CMD ["./bin/api"]

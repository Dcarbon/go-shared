# Images name: dcarbon/golang

FROM golang:1.23-alpine

RUN apk update && apk add --no-cache git ca-certificates alpine-sdk

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    echo "Load dependency successfully...!" 


WORKDIR /dcarbon
COPY . .

ENV PATH=$PATH:$GOPATH/bin
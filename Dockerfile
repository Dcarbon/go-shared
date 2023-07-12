# Images name: harbor.viet-tin.com/dcarbon/golang

FROM harbor.viet-tin.com/dcarbon/golang


WORKDIR /dcarbon/go-shared
COPY . .
RUN go mod tidy
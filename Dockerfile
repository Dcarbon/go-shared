# Images name: dcarbon/golang

FROM dcarbon/arch-proto:golang


WORKDIR /dcarbon/go-shared
COPY . .
RUN go mod tidy
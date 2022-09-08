FROM golang:1.19 AS build

WORKDIR /build

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

RUN make config

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build  -o nftasset

FROM scratch

COPY --from=build ["/build/config.yaml", "/"]
COPY --from=build ["/build/nft", "/"]

EXPOSE 8080

ENTRYPOINT ["/nft"]
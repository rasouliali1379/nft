FROM golang:1.19 AS build

WORKDIR /build

COPY . .

RUN mkdir docs
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.1
RUN swag fmt && swag init --parseDependency --parseInternal --parseDepth 1 -g app.go

RUN make config

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o nft ./cmd/app/main.go

FROM scratch

COPY --from=build ["/build/config.yaml", "/"]
COPY --from=build ["/build/nft", "/"]

EXPOSE 8080

ENTRYPOINT ["/nft"]
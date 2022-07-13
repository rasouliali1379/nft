FROM golang:1.18 AS build

WORKDIR /build


RUN mkdir -p /tmp

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

RUN make config

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build  -o maskan .

FROM scratch

COPY --from=build ["/build/config.yaml", "/"]
COPY --from=build ["/build/maskan", "/"]

EXPOSE 8080

ENTRYPOINT ["/maskan"]
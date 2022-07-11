FROM golang:1.18 AS build

WORKDIR /build


RUN mkdir -p /tmp

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build  -o maskan .
RUN make config

FROM scratch

COPY --from=build ["/build/config.yaml", "/"]
COPY --from=build ["/build/maskan", "/"]

# Declare volumes to mount
VOLUME /tmp

EXPOSE 8080

ENTRYPOINT ["/maskan"]
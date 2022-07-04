FROM golang:1.18 AS build

WORKDIR /build


RUN mkdir -p /tmp

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build  -o maskan .

FROM scratch

COPY --from=build ["/build/config.yaml", "/"]
COPY --from=build ["/build/maskan", "/"]

# Declare volumes to mount
VOLUME /tmp

EXPOSE 8080

ENTRYPOINT ["/maskan"]
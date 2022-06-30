FROM golang:1.18 AS builder

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

RUN go get -d -v

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/maskan

FROM scratch

COPY --from=builder /go/bin/maskan /go/bin/maskan

ENTRYPOINT ["/go/bin/hello"]
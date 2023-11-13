FROM golang:latest as builder

WORKDIR $GOPATH/src/app

COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app ./cmd/server/main.go

FROM alpine as production

COPY --from=builder /go/bin/app .

ENTRYPOINT ["./app"]

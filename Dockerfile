FROM golang:1.9-alpine AS builder

ENV BUILDROOT /go/src/github.com/adampointer/restservice
ADD . $BUILDROOT
WORKDIR $BUILDROOT

RUN go test -v ./...; \
    go build .

FROM alpine

ENV BUILDROOT /go/src/github.com/adampointer/restservice
COPY --from=builder $BUILDROOT/restservice /bin

RUN apk add -U ca-certificates

EXPOSE 8080

ENTRYPOINT ["/bin/restservice"]

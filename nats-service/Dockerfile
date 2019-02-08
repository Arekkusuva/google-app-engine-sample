FROM golang:1.11.5-alpine3.8 AS build
ENV GOBIN /go/bin
RUN apk add gcc g++ make git

RUN go get github.com/nats-io/gnatsd

FROM alpine:3.8
RUN apk add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/bin .

CMD ./gnatsd --port 8080

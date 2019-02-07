FROM golang:1.11.5-alpine3.8 AS build
ENV GOBIN /go/bin
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/google-app-egine-sample

# Copy source files
COPY go.mod .
COPY main.go .

# TODO: Run tests

# Compile source
RUN go build -o $GOBIN/service .

# Copy compiled app to run
FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/bin .

EXPOSE 8080
CMD ./service

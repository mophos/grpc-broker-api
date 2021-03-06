# Stage 1
FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go get -d -v
RUN go get github.com/moph-gateway/his-proto/proto
RUN go build -o rest-grpc .
# Stage 2
FROM alpine
RUN adduser -S -D -H -h /app siteslave
USER siteslave
COPY --from=builder /build/ /app/
WORKDIR /app
CMD ["./rest-grpc"]

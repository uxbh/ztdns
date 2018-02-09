FROM golang:1 AS build-env

WORKDIR /go/src/github.com/uxbh/ztdns
# Add source
COPY . .

# Install dependencies
RUN go get -d -v ./...
# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go install -v ./...

FROM alpine

# We need to add ca-certificates in order to make HTTPS API calls
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
# Copy binary
COPY --from=build-env /go/bin/ztdns .

ENTRYPOINT ["./ztdns", "--debug"]
CMD ["server"]
EXPOSE 53/udp

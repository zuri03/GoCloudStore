FROM golang:1.17 as builder
WORKDIR /go/src/github.com/zuri03/GoCloudStore/

COPY . .

RUN CGO_ENABLED=0 go build -o  build/Server ./cmd/server/main.go
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/zuri03/GoCloudStore/
COPY --from=builder /go/src/github.com/zuri03/GoCloudStore/build/Server ./main
EXPOSE 8080
ENTRYPOINT [ "./main" ]
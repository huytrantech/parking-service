FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
#download dependencies
RUN go mod download
COPY . .
RUN cd api && go build -o /go/bin/main
#remove comment if want apply unit test
#RUN apk add build-base
#RUN go test -v
FROM alpine:3.14
# RUN apk add --update ca-certificates
# RUN apk add --no-cache tzdata && \
#   cp -f /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
#   apk del tzdata
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/main .
COPY environment.local.env /app/environment.local.env
#COPY region.json /app/region.json
#expose port
EXPOSE 1323
ENTRYPOINT ["./main"]
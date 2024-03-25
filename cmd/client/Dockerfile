FROM golang:1.22.0-alpine as builder
RUN apk update && apk add git
WORKDIR /go/src/lokesh-katari/code-realm/cmd/client
COPY . .
RUN go mod init lokesh-katari/code-realm/cmd/client
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app .

FROM alpine
WORKDIR /app
COPY --from=builder /go/bin/app /app/app
ENTRYPOINT ["/app/app"]

FROM golang:1.22.0-alpine as builder

EXPOSE 55051

RUN apk update && apk add git

WORKDIR /go/src/lokesh-katari/code-realm/cmd/auth
COPY . .
RUN go mod init lokesh-katari/code-realm/cmd/auth
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app .

FROM scratch
WORKDIR /bin
COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]

# Build stage 
FROM golang:1.22.0-alpine as builder

RUN apk add --no-cache git

WORKDIR /go/src/lokesh-katari/code-realm/cmd/codeexecutor
COPY . .
RUN go mod init lokesh-katari/code-realm/cmd/codeexecutor
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app .

# Final stage (modified for DinD setup via a shell script)
FROM docker:20.10.17-dind

COPY --from=builder /go/bin/app /app
RUN apk add --no-cache ca-certificates tzdata

# Copy the start-dockerd.sh script
COPY start-dockerd.sh /start-dockerd.sh
RUN chmod +x /start-dockerd.sh 
# Set the entrypoint to execute the script
ENTRYPOINT ["/bin/sh", "-c", "/start-dockerd.sh && /app"]

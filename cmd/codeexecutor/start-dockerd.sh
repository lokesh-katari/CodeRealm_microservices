#!/bin/sh

# Start Docker daemon with privileged mode
dockerd --host=unix:///var/run/docker.sock --host=tcp://0.0.0.0:2375 & sleep 8

# Pull custom images (replace with your required images)
docker pull lokeshkatari/python-env:latest

sleep 2

docker pull lokeshkatari/js-env:latest

sleep 2
docker pull lokeshkatari/java-env:latest
sleep 2
docker pull lokeshkatari/gcc-env:latest
sleep 2
docker pull lokeshkatari/csharp-env:latest
sleep 2
docker pull lokeshkatari/php-env:latest
sleep 2
docker pull lokeshkatari/rust-env:latest
sleep 2
docker pull lokeshkatari/ruby-env:latest
sleep 2
# docker pull lokeshkatari/swift-env:latest
sleep 2
docker pull lokeshkatari/golang-env:latest

docker info

# Execute the main application
exec /app

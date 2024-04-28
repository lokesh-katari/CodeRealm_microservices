#!/bin/sh

# Start Docker daemon with privileged mode
dockerd --host=unix:///var/run/docker.sock --host=tcp://0.0.0.0:2375 & sleep 8

# Pull custom images (replace with your required images)
docker pull lokeshkatari/python-env:latest

sleep 3

docker pull lokeshkatari/js-env:latest
docker info

# Execute the main application
exec /app

docker rmi lokeshkatari/coderealm-client -f
docker buildx build --platform=linux/amd64   . -t lokeshkatari/coderealm-client:latest
docker push lokeshkatari/coderealm-client:latest


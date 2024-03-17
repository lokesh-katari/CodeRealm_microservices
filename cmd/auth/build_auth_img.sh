docker rmi lokeshkatari/coderealm-auth -f
docker buildx build --platform=linux/amd64   . -t lokeshkatari/coderealm-auth:latest
docker push lokeshkatari/coderealm-auth:latest


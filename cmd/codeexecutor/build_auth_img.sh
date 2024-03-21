docker rmi lokeshkatari/coderealm-codeexec -f
docker buildx build  --platform=linux/amd64   . -t lokeshkatari/coderealm-codeexec:latest
# docker push lokeshkatari/coderealm-codeExec:latest


docker rmi lokeshkatari/coderealm-codeexec -f
docker buildx build  --platform=linux/amd64   . -t lokeshkatari/coderealm-codeexec:latest
docker push lokeshkatari/coderealm-codeexec:latest



# to run the image in a privileged mode
# docker run --privileged -p 50052:50052  lokeshkatari/coderealm-codeexec

docker rmi lokeshkatari/coderealm-auth-envoy:latest -f
docker build . -t lokeshkatari/coderealm-auth-envoy:latest
docker push lokeshkatari/coderealm-auth-envoy:latest
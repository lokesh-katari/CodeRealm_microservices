docker rmi lokeshkatari/coderealm-auth-envoy:latest -f
docker build . -t lokeshkatari/coderealm-auth-envoy:latest
docker push lokeshkatari/coderealm-auth-envoy:latest
# docker run  -p 8000:8000 -p 9901:9901 --network test  lokeshkatari/coderealm-auth-envoy:latest
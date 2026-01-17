# Gateway (Client)

## Docker

1. Must use command for building image, proto needs to be compiled locally
```shell
docker buildx build --platform linux/amd64 --load -t gateway .
```
2. Run the image
```shell
docker run -d -p 5237:5236 --name gateway -e DATAMANAGER_HOST=datamanager:8080 gateway
```
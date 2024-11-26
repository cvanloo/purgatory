# Purgatory

## Docker Build

```sh
docker build -t go-purgatory .

docker run -it --rm --name purgatory \
    -e LISTEN_ADDR=:8000 \
    -e OUTSIDE_DOMAIN=https://purgatory.vanloo.ch/ \
    go-purgatory:latest
```

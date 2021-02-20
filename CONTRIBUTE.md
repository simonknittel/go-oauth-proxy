# How to contribute

## Build container image

```sh
docker build -t go-oauth-proxy:latest .
```

## Publish to Docker Hub

Link: <https://hub.docker.com/r/simonknittel/go-oauth-proxy>

```sh
docker tag go-oauth-proxy:latest simonknittel/go-oauth-proxy:latest
docker push simonknittel/go-oauth-proxy:latest
```

## Publish to Amazon ECR (Public)

Link: <https://gallery.ecr.aws/d2i9h0g7/simonknittel/go-oauth-proxy>

```sh
docker tag go-oauth-proxy:latest public.ecr.aws/d2i9h0g7/simonknittel/go-oauth-proxy:latest
docker push public.ecr.aws/d2i9h0g7/simonknittel/go-oauth-proxy:latest
```

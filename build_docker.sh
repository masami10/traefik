#!/bin/bash

version="0.0.1"

docker_repo="linshenqi/traefik"

docker build -t ${docker_repo}:$1 -t .

docker push ${docker_repo}:$1
docker push ${docker_repo}:latest

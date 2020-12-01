#!/bin/bash

echo "Compiling the Microservice.."
echo "This will use the golang:1.14 image from docker"
docker pull golang:1.14

echo "Compiling the API"
docker run -it --rm -v "$(pwd)":/go -e GOPATH= golang:1.14 sh -c "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -a --installsuffix cgo --ldflags='-s'"
mv ./go-chat ./docker/
docker build -t go-chat:latest docker/

echo "Copying config json"
cp ./config/config.json.example ./config/config.json

echo "Checking Swarm network"
if ! docker node ls &>/dev/null ; then
    docker swarm init
fi
if ! docker network inspect main &>/dev/null ; then
    docker network create -d overlay main
fi

docker stack deploy -c build/docker-compose.yaml gcl

echo "It's all set, ready to go!"
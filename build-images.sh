#!/bin/bash

# Remove all Docker images
sudo docker rmi -f $(sudo docker images -a -q)

# Variables for the Docker images
GO_CLIENT_IMAGE="go-client-grpc"
RUST_CLIENT_IMAGE="rust-client-grpc"


GO_SERVER_NATACION_IMAGE="golang-server-natacion-grpc"

GO_SERVER_WINNERS_IMAGE="consumer-winners"
GO_SERVER_LOSERS_IMAGE="consumer-losers"
GO_SERVER_ATLETISMO_IMAGE="golang-server-atletismo-grpc"
GO_SERVER_BOXEO_IMAGE="golang-server-boxeo-grpc"

DOCKERHUB_USERNAME="mrpony21"
TAG="0.2"

# # # Build the Docker image for the Go client
# sudo docker build -t $GO_CLIENT_IMAGE ./gRPC/go-client
# # Build the Docker image for the Rust client
# sudo docker build -t $RUST_CLIENT_IMAGE ./gRPC/grpc-client

# Build the Docker image for the Go server Natacion
#sudo docker build -t $GO_SERVER_NATACION_IMAGE ./gRPC/go-server-natacion
sudo docker build -t $GO_SERVER_WINNERS_IMAGE ./consumer/consumer-winners
sudo docker build -t $GO_SERVER_LOSERS_IMAGE ./consumer/consumer-losers
# Build the Docker image for the Go server Atletismo
# sudo docker build -t $GO_SERVER_ATLETISMO_IMAGE ./gRPC/go-server-atletismo
# # # Build the Docker image for the Go server Boxeo
# sudo docker build -t $GO_SERVER_BOXEO_IMAGE ./gRPC/go-server-boxeo



# Tag the Docker image
# docker tag $GO_CLIENT_IMAGE "$DOCKERHUB_USERNAME/$GO_CLIENT_IMAGE:$TAG"
# docker tag $RUST_CLIENT_IMAGE "$DOCKERHUB_USERNAME/$RUST_CLIENT_IMAGE:$TAG"

docker tag $GO_SERVER_WINNERS_IMAGE "$DOCKERHUB_USERNAME/$GO_SERVER_WINNERS_IMAGE:$TAG"
docker tag $GO_SERVER_LOSERS_IMAGE "$DOCKERHUB_USERNAME/$GO_SERVER_LOSERS_IMAGE:$TAG"
# docker tag $GO_SERVER_NATACION_IMAGE "$DOCKERHUB_USERNAME/$GO_SERVER_NATACION_IMAGE:$TAG"
# docker tag $GO_SERVER_ATLETISMO_IMAGE "$DOCKERHUB_USERNAME/$GO_SERVER_ATLETISMO_IMAGE:$TAG"
# docker tag $GO_SERVER_BOXEO_IMAGE "$DOCKERHUB_USERNAME/$GO_SERVER_BOXEO_IMAGE:$TAG"

# # Push the Docker image to DockerHub
# docker push "$DOCKERHUB_USERNAME/$GO_CLIENT_IMAGE:$TAG"
# docker push "$DOCKERHUB_USERNAME/$RUST_CLIENT_IMAGE:$TAG"

docker push "$DOCKERHUB_USERNAME/$GO_SERVER_WINNERS_IMAGE:$TAG"
docker push "$DOCKERHUB_USERNAME/$GO_SERVER_LOSERS_IMAGE:$TAG"
# docker push  "$DOCKERHUB_USERNAME/$GO_SERVER_NATACION_IMAGE:$TAG"
# docker push "$DOCKERHUB_USERNAME/$GO_SERVER_ATLETISMO_IMAGE:$TAG"
# docker push "$DOCKERHUB_USERNAME/$GO_SERVER_BOXEO_IMAGE:$TAG"

echo "Docker images pushed successfully."
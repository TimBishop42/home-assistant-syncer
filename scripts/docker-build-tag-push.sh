#!/bin/bash

# Ensure the script runs from the root of the project directory
cd ../

# Set image names
IMAGE_NAME="homeassistant/syncer"
DOCKERHUB_REPO="tbished/syncer"
TAG="latest"
DOCKERFILE_PATH="docker/Dockerfile"

# Build the Docker image
echo "Building Docker image..."
docker build -t $IMAGE_NAME -f $DOCKERFILE_PATH .
if [ $? -ne 0 ]; then
  echo "Docker build failed."
  exit 1
fi

# Tag the Docker image
echo "Tagging Docker image..."
docker tag $IMAGE_NAME $DOCKERHUB_REPO:$TAG
if [ $? -ne 0 ]; then
  echo "Docker tag failed."
  exit 1
fi

# Push the Docker image to Docker Hub
echo "Pushing Docker image to Docker Hub..."
docker push $DOCKERHUB_REPO:$TAG
if [ $? -ne 0 ]; then
  echo "Docker push failed."
  exit 1
fi

echo "Docker image built, tagged, and pushed successfully!"
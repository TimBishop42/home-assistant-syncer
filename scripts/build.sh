#!/bin/bash

# Build the Docker image
docker build -t my-go-service .

# Run the Docker container
docker run --rm -it my-go-service

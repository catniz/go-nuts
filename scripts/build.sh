#!/bin/bash

# ex) scripts/build.sh bullish 1.0.0

source "$(dirname "$0")/make_env.sh"

TARGET=$1
TAG=$2

dockerfile=$([ -f "cmd/$TARGET/Dockerfile" ] && echo "cmd/$TARGET/Dockerfile" || echo "build/common.Dockerfile")

echo "Building $TARGET image with $dockerfile"

docker build -t "$ECR_USERNAME/go-nuts/$TARGET:$TAG" -f "$dockerfile" . --build-arg TARGET="$TARGET"
#docker push "$ECR_USERNAME/$TARGET:$TAG"

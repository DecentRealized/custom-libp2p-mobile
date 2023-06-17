#!/usr/bin/env sh

GOOS=${GOOS:-linux}
GOARCH=${GOARCH:-amd64}
GOVERSION=${GOVERSION:-1.17.1}
TAG=${TAG:-main}
CMD=${CMD:-make build}
ARGS=${ARGS:-}
DOCKER_IMAGE_CROSS=docker.elastic.co/beats-dev/golang-crossbuild:${GOVERSION}

podman run -t --rm -v "$PWD":/app -w /app \
	-e CGO_ENABLED=1 "${ARGS}" \
	"${DOCKER_IMAGE_CROSS}-${TAG}" \
	--build-cmd "${CMD}" -p "${GOOS}/${GOARCH}"
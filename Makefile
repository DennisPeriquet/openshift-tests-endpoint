VER := $(shell cat VER)
QUAY_USER := $(shell whoami)

all: build

build:
	go build ./cmd/endpoint_server/

podman-build:
	podman build -t quay.io/${QUAY_USER}/openshift_tests_endpoint:${VER} .

podman-tag:
	podman tag quay.io/${QUAY_USER}/openshift_tests_endpoint:${VER} docker.io/${QUAY_USER}/openshift_tests_endpoint:${VER}

docker-push:
	podman push docker.io/${QUAY_USER}/openshift_tests_endpoint:${VER}

clean:
	rm -f endpoint_server

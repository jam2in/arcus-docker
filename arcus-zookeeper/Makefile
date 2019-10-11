include Docker.env
export $(shell sed 's/=.*//' Docker.env)

ARCUS_ZOOKEEPER_IMAGE=${ARCUS_ZOOKEEPER_IMAGE_REPO}:${ARCUS_ZOOKEEPER_IMAGE_TAG}

build:
	@ sed 's|\$${ARCUS_ZOOKEEPER_VERSION}|${ARCUS_ZOOKEEPER_VERSION}|g' Dockerfile.template > Dockerfile
	@ docker build -t ${ARCUS_ZOOKEEPER_IMAGE} . --network=host

push:
	@ docker push ${ARCUS_ZOOKEEPER_IMAGE}

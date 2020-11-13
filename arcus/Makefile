include Docker.env
export $(shell sed 's/=.*//' Docker.env)

ARCUS_IMAGE=${ARCUS_IMAGE_REPO}:${ARCUS_IMAGE_TAG}

build:
	@ sed 's|\$${ARCUS_MEMCACHED_VERSION}|${ARCUS_MEMCACHED_VERSION}|g' Dockerfile.template > Dockerfile
	@ docker build -t ${ARCUS_IMAGE} . --network=host

push:
	@ docker push ${ARCUS_IMAGE}

include Docker.env
export $(shell sed 's/=.*//' Docker.env)

ARCUS_MEMCACHED_IMAGE=${ARCUS_MEMCACHED_IMAGE_REPO}:${ARCUS_MEMCACHED_IMAGE_TAG}

build:
	@ sed 's|\$${ARCUS_MEMCACHED_VERSION}|${ARCUS_MEMCACHED_VERSION}|g' Dockerfile.template > Dockerfile
	@ docker build -t ${ARCUS_MEMCACHED_IMAGE} . --network=host

push:
	@ docker push ${ARCUS_MEMCACHED_IMAGE}

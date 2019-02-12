IMAGE=jam2in/arcus-zookeeper:3.4.6

build:
	# alpine Repository Issue
	# https://github.com/gliderlabs/docker-alpine/issues/334
	docker build -t ${IMAGE} . --network=host

push:
	docker push ${IMAGE}

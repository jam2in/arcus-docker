IMAGE=arcus-zookeeper:latest
CONTAINER=arcus-zookeeper

build:
	docker build -t ${IMAGE} .

rund:
	docker run --name ${CONTAINER} --env-file=./env -d ${IMAGE}

runit:
	docker run --name ${CONTAINER} --env-file=./env -it ${IMAGE} bash

exec:
	docker exec -it ${CONTAINER} bash

rmi:
	docker rmi ${IMAGE}

rm:
	docker rm ${CONTAINER} -f

# arcus-docker

## How to create arcus-ubuntu:14.04 docker image
```
$ cd arcus-ubuntu-14.04
$ make build
```

## How to create arcus-zookeeper docker images
```
$ cd arcus-zookeeper
$ make build
```

## How to apply kubernetes statefulset of arcus-zookeeper
```
$ cd arcus-kubernetes
$ kubectl apply -f arcus-zookeeper.yaml
```

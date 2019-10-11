# arcus-docker

## Building arcus-memcached docker image from Docker.env and Makefile

1. Edit Docker.env file
```
$ cd arcus-memcached

$ cat Docker.env
ARCUS_MEMCACHED_VERSION=1.11.7
ARCUS_MEMCACHED_IMAGE_REPO=jam2in/arcus-memcached
ARCUS_MEMCACHED_IMAGE_TAG=1.11.7
```

2. Build docker image
```
$ cd arcus-memcached

$ make build
```

3. Push docker image
```
$ docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username (xxx): xxx
Password: xxx

$ cd arcus-memcached

$ make push
```

## Building arcus-zookeeper docker image from Docker.env and Makefile

1. Edit Docker.env file
```
$ cd arcus-zookeeper

$ cat Docker.env
ARCUS_ZOOKEEPER_VERSION=3.4.7
ARCUS_ZOOKEEPER_IMAGE_REPO=jam2in/arcus-zookeeper
ARCUS_ZOOKEEPER_IMAGE_TAG=3.4.7
```

2. Build docker image
```
$ cd arcus-zookeeper

$ make build
```

3. Push docker image
```
$ docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username (xxx): xxx
Password: xxx

$ cd arcus-zookeeper

$ make push
```

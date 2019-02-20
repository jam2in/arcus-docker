FROM centos:6.9

ARG ARCUS_REPOSITORY=https://github.com/naver/arcus.git

RUN yum update -y; \
    yum install -y \
      gcc \
      gcc-c++ \
      autoconf \
      automake \
      libtool \
      pkgconfig \
      cppunit-devel \
      python-setuptools \
      python-devel \
      java-1.8.0-openjdk-devel \
      ant \
      git \
      nc

ENV ARCUS_USER    arcus
ENV ARCUS_DIR     /arcus
ENV MEMCACHED_DIR /arcus/server
ENV ZOOKEEPER_DIR /arcus/zookeeper
ENV PATH          $PATH:$ARCUS_DIR/server:$ZOOKEEPER_DIR/bin

RUN set -ex; \
    adduser "$ARCUS_USER"; \
    git clone $ARCUS_REPOSITORY; \
    $ARCUS_DIR/scripts/build.sh; \
    chown -R "$ARCUS_USER:$ARCUS_USER" "$ARCUS_DIR";

WORKDIR $ARCUS_DIR

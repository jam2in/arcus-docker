FROM centos:6.9

ARG ARCUS_REPOSITORY=https://github.com/naver/arcus.git
ARG MEMCACHED_REPOSITORY=https://github.com/naver/arcus-memcached.git

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
ENV ZOOKEEPER_DIR /arcus/zookeeper
ENV MEMCACHED_DIR /arcus-memcached
ENV PATH          $PATH:$ZOOKEEPER_DIR/bin:$MEMCACHED_DIR

RUN set -ex; \
    adduser "$ARCUS_USER"; \
    git clone $ARCUS_REPOSITORY; \
    $ARCUS_DIR/scripts/build.sh; \
    chown -R "$ARCUS_USER:$ARCUS_USER" "$ARCUS_DIR"

RUN set -ex; \
    git clone $MEMCACHED_REPOSITORY; \
    pushd $MEMCACHED_DIR; \
    git checkout 1.11.6; \
    config/autorun.sh; \
    ./configure --enable-zk-integration --with-libevent=$ARCUS_DIR --with-zookeeper=$ARCUS_DIR; \
    make; \
    popd; \
    chown -R "$ARCUS_USER:$ARCUS_USER" "$MEMCACHED_DIR"

# image
FROM jam2in/arcus-ubuntu:14.04

# maintainer
MAINTAINER hjyun <hjyun@jam2in.com>

# env
ENV ARCUS_PATH        /usr/lib/arcus
ENV ARCUS_USER        arcus
ENV ZOOKEEPER_PATH    $ARCUS_PATH/zookeeper

# add arcus user
RUN useradd -d $ARCUS_PATH -c ARCUS -s /bin/bash $ARCUS_USER
RUN rm -rf $ARCUS_PATH/*
RUN echo "$ARCUS_USER:$ARCUS_USER" : chpasswd

# download arcus
WORKDIR /usr/lib
RUN git clone https://github.com/hjyun328/arcus.git
RUN chown -R $ARCUS_USER:$ARCUS_USER $ARCUS_PATH

# build arcus
USER $ARCUS_USER
WORKDIR $ARCUS_PATH/scripts
RUN ./build.sh
WORKDIR $ARCUS_PATH

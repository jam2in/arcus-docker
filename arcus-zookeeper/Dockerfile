FROM openjdk:8-jre-alpine

ARG ZOOKEEPER_REPOSITORY=http://archive.apache.org/dist/zookeeper/zookeeper-3.4.7/zookeeper-3.4.7.tar.gz
ARG ZOOKEEPER_FILE=zookeeper-3.4.7.tar.gz

RUN apk add --no-cache \
    bash \
    su-exec

ENV ZOOKEEPER_USER     zookeeper
ENV ZOOKEEPER_DIR      /zookeeper-3.4.7
ENV ZOOKEEPER_CONF_DIR $ZOOKEEPER_DIR/conf
ENV ZOOKEEPER_DATA_DIR $ZOOKEEPER_DIR/data
ENV ZOOCFGDIR          $ZOOKEEPER_CONF_DIR
ENV PATH               $PATH:$ZOOKEEPER_DIR/bin

RUN set -ex; \
    adduser -D "$ZOOKEEPER_USER"; \
    wget -q "$ZOOKEEPER_REPOSITORY"; \
    tar -xzf "$ZOOKEEPER_FILE"; \
    mkdir -p "$ZOOKEEPER_CONF_DIR" "$ZOOKEEPER_DATA_DIR"; \
    chown -R "$ZOOKEEPER_USER:$ZOOKEEPER_USER" "$ZOOKEEPER_DIR" "$ZOOKEEPER_CONF_DIR" "$ZOOKEEPER_DATA_DIR"; \
    rm -f "$ZOOKEEPER_FILE";

WORKDIR $ZOOKEEPER_DIR

FROM openjdk:8-jre-alpine

RUN apk add --no-cache \
    bash \
    su-exec

ENV ZOOKEEPER_USER=zookeeper \
    ZOOKEEPER_DIR=/zookeeper-3.4.7 \
    ZOOKEEPER_CONF_DIR=/conf \
    ZOOKEEPER_DATA_DIR=/data

ENV ZOOCFGDIR $ZOOKEEPER_CONF_DIR

RUN set -ex; \
    adduser -D "$ZOOKEEPER_USER"; \
    mkdir -p "$ZOOKEEPER_CONF_DIR" "$ZOOKEEPER_DATA_DIR";

RUN set -ex; \
    wget -q "http://archive.apache.org/dist/zookeeper/zookeeper-3.4.7/zookeeper-3.4.7.tar.gz"; \
    tar -xzf "zookeeper-3.4.7.tar.gz"; \
    chown -R "$ZOOKEEPER_USER:$ZOOKEEPER_USER" "$ZOOKEEPER_DIR"; \
    cp "$ZOOKEEPER_DIR/conf/"* "$ZOOKEEPER_CONF_DIR"; \
    rm -f "zookeeper-3.4.7.tar.gz";

WORKDIR $ZOOKEEPER_DIR

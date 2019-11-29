FROM centos:7

ENV ARCUS_APPLICATION_SAMPLE_USER     arcus
ENV ARCUS_APPLICATION_SAMPLE_DIR      /arcus-application-sample

COPY arcus-application-sample $ARCUS_APPLICATION_SAMPLE_DIR

RUN yum update -y; \
    yum install -y \
        java-1.8.0-openjdk-devel \
        curl

RUN set -ex; \
    adduser "$ARCUS_APPLICATION_SAMPLE_USER"; \
    chown -R "$ARCUS_APPLICATION_SAMPLE_USER:$ARCUS_APPLICATION_SAMPLE_USER" "$ARCUS_APPLICATION_SAMPLE_DIR"

WORKDIR $ARCUS_APPLICATION_SAMPLE_DIR

CMD ["java", "-jar", "war/arcus-application-sample-0.0.1.war"]

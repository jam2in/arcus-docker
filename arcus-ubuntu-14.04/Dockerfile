# image
FROM ubuntu:14.04

# maintainer
MAINTAINER hjyun <hjyun328@gmail.com>

# env 
ENV DEBIAN_FRONTEND   noninteractive
ENV JAVA_HOME         /usr/lib/jvm/java-8-openjdk-amd64

# install dependency
RUN apt-get update
RUN apt-get install -y build-essential autoconf automake libtool libcppunit-dev python-setuptools python-dev
RUN apt-get install -y software-properties-common
RUN add-apt-repository -y ppa:openjdk-r/ppa
RUN apt-get update
RUN apt-get install -y openjdk-8-jdk
RUN apt-get install -y ant
RUN apt-get install -y git

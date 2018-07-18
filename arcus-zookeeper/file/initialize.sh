#!/bin/bash

ZOOKEEPER_CONF=$ZOOKEEPER_PATH/conf/zoo.cfg
ZOOKEEPER_DATA=$ZOOKEEPER_PATH/data/myid

HOST_NAME_SHORT=`hostname -s`
HOST_NAME_DOMAIN=`hostname -d`
if [[ $HOST_NAME_SHORT =~ (.*)-([0-9]+)$ ]]; then
  HOST_NAME=${BASH_REMATCH[1]}
  HOST_ORG=${BASH_REMATCH[2]}
fi

function create_config() {
    mkdir -p $ZOOKEEPER_PATH/conf
    rm -f $ZOOKEEPER_CONF
    echo "maxClientCnxns=$ZOOKEEPER_CONF_MAX_CLIENT_CNXNS"       >> $ZOOKEEPER_CONF
    echo "tickTime=$ZOOKEEPER_CONF_TICK_TIME"                    >> $ZOOKEEPER_CONF
    echo "initLimit=$ZOOKEEPER_CONF_INIT_LIMIT"                  >> $ZOOKEEPER_CONF
    echo "syncLimit=$ZOOKEEPER_CONF_SYNC_LIMIT"                  >> $ZOOKEEPER_CONF
    echo "dataDir=$ZOOKEEPER_PATH/data"                          >> $ZOOKEEPER_CONF
    echo "clientPort=$ZOOKEEPER_CONF_CLIENT_PORT"                >> $ZOOKEEPER_CONF
    echo "minSessionTimeout=$ZOOKEEPER_CONF_MIN_SESSION_TIMEOUT" >> $ZOOKEEPER_CONF
    echo "maxSessionTimeout=$ZOOKEEPER_CONF_MAX_SESSION_TIMEOUT" >> $ZOOKEEPER_CONF

    for (( i=0; i<$ZOOKEEPER_CONF_SERVERS; i++ ))
    do
        echo "server.$((i+1))=$HOST_NAME-$i.$HOST_NAME_DOMAIN:$ZOOKEEPER_CONF_SERVER_PORT:$ZOOKEEPER_CONF_LEADER_ELECTION_PORT" >> $ZOOKEEPER_CONF
    done
}

function create_data() {
    mkdir -p $ZOOKEEPER_PATH/data
    echo "$((HOST_ORG+1))" > $ZOOKEEPER_DATA
}

create_config && create_data

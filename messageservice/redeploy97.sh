#!/bin/bash

set -e

version_from_conf=`awk '$1=="versioncode" {print $3}' ./conf/app.conf`

if [ $# == 0 ] && [ -z $version_from_conf ]; then
    echo "baby,we need a version code"
    exit 1
fi

version=$version_from_conf
if [ $# == 1 ]; then
    version=$1
fi

echo $version

default_runmode="test"
runmode=`awk '$1=="runmode" {print $3}' ./conf/app.conf`

if [ $default_runmode != $runmode ]
then
    echo "$runmode is err,you should in $default_runmode"
	exit 1
fi

echo "redeploy messageservice"

cd GoWorkTest/src/zhangmai_micro/messageservice/
if docker build -t messageservicetest:$version .
then
    echo "stop and rm old container,start new one..."
    docker stop messageservicetest
    docker rm messageservicetest
    docker run --restart=always --name  messageservicetest -v /opt/lvdoutest/zhangmailog/messageservicelog:/opt/zhangmaiLogs -v /data:/data -d -p 9010:9010 -p 9011:9011 messageservicetest:$version
    docker ps
    docker rmi -f  `docker images | grep '<none>' | awk '{print $3}'`
fi



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

echo "redeploy taskservicetest"

cd GoWorkTest/src/zhangmai_micro/taskservice/
if docker build -t taskservicetest:$version .
then
    echo "stop and rm old container,start new one..."
    docker stop taskservicetest
    docker rm taskservicetest
    docker run --restart=always --name  taskservicetest -v /opt/lvdoutest/zhangmailog/taskservicelog:/opt/zhangmaiLogs -v /data:/data -d -p 9016:9016 -p 9017:9017 taskservicetest:$version
    docker ps
    docker rmi -f  `docker images | grep '<none>' | awk '{print $3}'`
fi



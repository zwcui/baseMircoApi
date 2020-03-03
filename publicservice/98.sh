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

default_runmode="prod"
runmode=`awk '$1=="runmode" {print $3}' ./conf/app.conf`

if [ $default_runmode != $runmode ]
then
    echo "$runmode is err,you should in $default_runmode"
	exit 1
fi

go clean

echo "rsync publicservice"
rsync -avzIP --delete --exclude .git ./  root@#ip#:~/GoWork/src/zhangmai_micro/publicservice/

ssh  root@#ip# version=$version 'bash -se' <<'ENDSSH'
cd GoWork/src/zhangmai_micro/publicservice/
if docker build -t jingting_publicservice:$version .
then
    echo "stop and rm old container,start new one..."
    docker stop jingting_publicservice
    docker rm jingting_publicservice
    docker run --restart=always --name jingting_publicservice -v /opt/91lvdou/zhangmaimicrolog/publicservicelog:/opt/zhangmaiLogs -v /data:/data -d -p 9012:9012 -p 9013:9013 jingting_publicservice:$version
    docker ps
    docker rmi -f  `docker images | grep '<none>' | awk '{print $3}'`
fi
ENDSSH



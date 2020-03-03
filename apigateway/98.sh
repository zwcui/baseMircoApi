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

echo "rsync apigateway"
rsync -avzIP --delete --exclude .git ./  root@#ip#:~/GoWork/src/zhangmai_micro/apigateway/

ssh  root@#ip# version=$version 'bash -se' <<'ENDSSH'
cd GoWork/src/zhangmai_micro/apigateway/
if docker build -t jingting_apigateway:$version .
then
    echo "stop and rm old container,start new one..."
    docker stop jingting_apigateway
    docker rm jingting_apigateway
    docker run --restart=always --name jingting_apigateway -d -p 9999:9999 -p 9030:9030 -p 9021:9021 jingting_apigateway:$version
    docker ps
    docker rmi -f  `docker images | grep '<none>' | awk '{print $3}'`
fi
ENDSSH



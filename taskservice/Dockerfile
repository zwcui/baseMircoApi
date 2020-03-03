FROM golang:1.12.9-alpine

MAINTAINER zwcui<zwcui2017@163.com>

ENV kpdir /go/src/jingting_server/taskservice

RUN mkdir -p ${kpdir}

ADD . ${kpdir}/

WORKDIR ${kpdir}

RUN go build -v

EXPOSE 9016

ENTRYPOINT ["./taskservice"]
FROM golang:1.12.9-alpine

MAINTAINER zwcui<zwcui2017@163.com>

ENV kpdir /go/src/jingting_server/socketservice

RUN mkdir -p ${kpdir}

ADD . ${kpdir}/

WORKDIR ${kpdir}

RUN go build -v

EXPOSE 9019

ENTRYPOINT ["./socketservice"]
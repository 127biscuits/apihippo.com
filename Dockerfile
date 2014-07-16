FROM ubuntu:14.04
MAINTAINER Ivan Pedrazas <ivan@pedrazas.me>


RUN apt-get update

RUN apt-get install -y curl git bzr mercurial golang

RUN mkdir -p /usr/local/go

ENV GOPATH  /go
ENV GOROOT  /usr/local/go


WORKDIR /go/src/127biscuits/apihippo.com
RUN go get github.com/127biscuits/apihippo.com
ADD . /go/src/github.com/127biscuits/apihippo.com
RUN go get
RUN go build

EXPOSE 8000

ENTRYPOINT ./apihippo.com -option value args
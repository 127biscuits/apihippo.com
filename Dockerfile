FROM ubuntu:14.04
MAINTAINER Ivan Pedrazas <ivan@pedrazas.me>


RUN apt-get update

RUN apt-get install -y curl git bzr mercurial


RUN curl -s https://storage.googleapis.com/golang/go1.3.linux-amd64.tar.gz | tar  -v -C /usr/local/ -xz


ENV PATH  /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
ENV GOPATH  /go
ENV GOROOT  /usr/local/go


WORKDIR /go/src/127biscuits/apihippo.com
RUN go get github.com/127biscuits/apihippo.com
ADD . /go/src/github.com/127biscuits/apihippo.com
RUN go get
RUN go build

EXPOSE 8000

ENTRYPOINT ./apihippo.com -option value args
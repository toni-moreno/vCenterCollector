FROM alpine:latest
MAINTAINER Toni Moreno <toni.moreno@gmail.com>

ADD ./vcentercollector-last.tar.gz /

VOLUME ["/opt/vcentercollector/conf", "/opt/vcentercollector/log"]

EXPOSE 8090

WORKDIR /opt/vcentercollector
COPY ./vcentercollector.toml ./conf/
COPY ./start.sh /

ENTRYPOINT ["/start.sh"]

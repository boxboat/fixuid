FROM alpine:latest

RUN addgroup -g 1000 docker && \
    adduser -u 1000 -G docker -h /home/docker -s /bin/sh -D docker

COPY stage /
RUN chmod u+s /usr/local/bin/fixuid && \
    chown -R docker:docker /tmp/*

USER docker:docker

RUN touch /home/docker/aaa && \
    touch /home/docker/zzz

FROM centos:latest

RUN groupadd -g 1000 docker && \
    useradd -u 1000 -g docker -d /home/docker -s /bin/sh docker

COPY stage /
RUN chmod u+s /usr/local/bin/fixuid && \
    chown -R docker:docker /tmp/*

USER docker:docker

RUN touch /home/docker/aaa && \
    touch /home/docker/zzz

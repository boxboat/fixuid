FROM debian:latest

RUN addgroup --gid 1000 docker && \
    adduser --uid 1000 --ingroup docker --home /home/docker --shell /bin/sh --disabled-password --gecos "" docker

COPY stage /
RUN chmod u+s /usr/local/bin/fixuid && \
    chown -R docker:docker /tmp/*

USER docker:docker

RUN touch /home/docker/aaa && \
    touch /home/docker/zzz

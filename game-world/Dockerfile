FROM alpine:latest
MAINTAINER  jzy
WORKDIR /go/src/game-world
RUN mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY WorldApp /go/src/game-world/WorldApp
COPY config/*.json /go/src/game-world/config/

CMD /go/src/game-world/WorldApp ${GO_OPTS}
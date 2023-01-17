FROM alpine:latest
MAINTAINER  jzy
WORKDIR /go/src/game-gate
RUN mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2


COPY GateApp /go/src/game-gate/GateApp
COPY config/*.json /go/src/game-gate/config/


CMD /go/src/game-gate/GateApp ${GO_OPTS}
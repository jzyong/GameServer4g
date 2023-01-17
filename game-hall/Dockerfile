FROM alpine:latest
MAINTAINER  jzy
WORKDIR /go/src/game-hall
RUN mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY HallApp /go/src/game-hall/HallApp
COPY config/*.json /go/src/game-hall/config/


#CMD ["/bin/bash", "/go/src/kdnn-gate-go/src/script/run.sh"]
CMD /go/src/game-hall/HallApp ${GO_OPTS}
FROM golang
MAINTAINER  jzy
WORKDIR /go/src/
COPY . game-world
EXPOSE 8889


#CMD ["/bin/bash", "/go/src/kdnn-gate-go/src/script/run.sh"]
CMD game-world/WorldApp ${GO_OPTS}
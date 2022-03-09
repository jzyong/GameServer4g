FROM golang
MAINTAINER  jzy
WORKDIR /go/src/
COPY . world
EXPOSE 8889


#CMD ["/bin/bash", "/go/src/kdnn-gate-go/src/script/run.sh"]
CMD world/WorldApp ${GO_OPTS}
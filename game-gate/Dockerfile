FROM golang
MAINTAINER  jzy
WORKDIR /go/src/
COPY . gate
EXPOSE 8888


#CMD ["/bin/bash", "/go/src/kdnn-gate-go/src/script/run.sh"]
CMD gate/GateApp ${GO_OPTS}
FROM golang
MAINTAINER  jzy
WORKDIR /go/src/
COPY . hall
EXPOSE 8888


#CMD ["/bin/bash", "/go/src/kdnn-gate-go/src/script/run.sh"]
CMD hall/HallApp ${GO_OPTS}
call build_linux.bat
docker image build -t mmo-gate .

set goRunParams="-config /go/src/gate/config/GateConfig.json"

docker stop mmo-gate
docker rm mmo-gate
docker run --name mmo-gate -e GO_OPTS=%goRunParams% mmo-gate

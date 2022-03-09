call build_linux.bat
docker image build -t mmo-hall .

set goRunParams="-config /go/src/hall/config/HallConfig.json"

docker stop mmo-hall
docker rm mmo-hall
docker run --name mmo-hall -e GO_OPTS=%goRunParams% mmo-hall

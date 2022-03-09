call build_linux.bat
docker image build -t mmo-world .

set goRunParams="-config /go/src/world/config/WorldConfig.json"

docker stop mmo-world
docker rm mmo-world
docker run --name mmo-world -e GO_OPTS=%goRunParams% mmo-world

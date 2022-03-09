call build_linux.bat
docker image build -t game-world .

set goRunParams="-config /go/src/game-world/config/ApplicationConfig_develop.json"

docker stop game-world
docker rm game-world
docker run --name game-world -e GO_OPTS=%goRunParams% game-world

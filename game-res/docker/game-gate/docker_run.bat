call build_linux.bat
docker image build -t game-gate .

set goRunParams="-config /go/src/game-gate/config/ApplicationConfig_develop.json"

docker stop game-gate
docker rm game-gate
docker run --name game-gate -e GO_OPTS=%goRunParams% game-gate

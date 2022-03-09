call build_linux.bat
docker image build -t game-hall .

set goRunParams="-config /go/src/game-hall/config/ApplicationConfig_develop.json"

docker stop game-hall
docker rm game-hall
docker run --name game-hall -e GO_OPTS=%goRunParams% game-hall

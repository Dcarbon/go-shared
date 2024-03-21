cd ../arch-proto
./scripts/docker-go.sh
cd -

docker build -t dcarbon/go-shared .
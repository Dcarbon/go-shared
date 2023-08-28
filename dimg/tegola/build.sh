TAG=harbor.viet-tin.com/dcarbon/tegola
docker build -t $TAG .

if [[ "$1" == "push" ]]; then
    docker push $TAG
fi
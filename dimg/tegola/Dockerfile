FROM dcarbon/golang as builder

WORKDIR /tmp

RUN git clone https://github.com/go-spatial/tegola.git &&  \
    cd tegola/cmd/tegola && \
    go build -buildvcs=false && \
    echo "Build tegola done ...!"

FROM dcarbon/dimg:minimal
COPY --from=builder /tmp/tegola/cmd/tegola/tegola /usr/bin/tegola

COPY ./run_tegola.sh /usr/bin/
COPY ./tegola_config.toml /etc/conf.d/

ENV REDIS_SSL=false
ENV REDIS_PASSWORD=""
ENV TEGOLA_CONFIG=/etc/conf.d/tegola_config.toml

ENTRYPOINT [ "sh", "-c", "run_tegola.sh" ]
FROM tendermint/tendermint:latest
USER root
COPY ./init/docker-entrypoint.sh /usr/local/bin/
COPY ./init/node_key.json config/node_key.json
RUN chmod u+x /usr/local/bin/docker-entrypoint.sh

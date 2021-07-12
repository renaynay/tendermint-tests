FROM tendermint/tendermint:latest
USER root
COPY ./init/node_key.json ~/config/node_key.json
COPY ./init/docker-entrypoint.sh /usr/local/bin/
RUN chmod u+x /usr/local/bin/docker-entrypoint.sh

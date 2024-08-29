FROM nemo-network-base

COPY ./testing/e2etest-local/local.sh /nemo-network/local.sh
COPY ./testing/genesis.sh /nemo-network/genesis.sh
COPY ./testing/start.sh /nemo-network/start.sh
COPY ./daemons/pricefeed/client/constants/testdata /nemo-network/exchange_config
COPY ./testing/delaymsg_config /nemo-network/delaymsg_config

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /nemo-network/local.sh

ENV DAEMON_NAME=nemod
ENTRYPOINT ["cosmovisor", "run"]

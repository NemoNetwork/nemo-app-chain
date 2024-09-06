FROM nemo-network-base

RUN apk add --no-cache bash jq aws-cli
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

COPY ./testing/version/. /nemo-network/
COPY ./testing/mainnet/. /nemo-network/

ENV HOME /nemo-network
WORKDIR $HOME

RUN /nemo-network/mainnet.sh

ENTRYPOINT ["/nemo-network/start.sh"]

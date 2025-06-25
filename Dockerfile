FROM ubuntu:20.04 AS build-env

# Set up dependencies
ENV PACKAGES="curl wget make git build-essential libc6-dev gcc g++ \
    libudev-dev python3 ca-certificates libgmp-dev flex bison"

# Install minimum necessary dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends $PACKAGES \
    && update-ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install pdc
RUN wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz \
    && tar -xf pbc-0.5.14.tar.gz \
    && cd pbc-0.5.14/ \
    && ./configure \
    && make \
    && make install \
    && ldconfig / \
    && cd .. && rm -rf pbc-0.5.14/ pbc-0.5.14.tar.gz

# Install Go 1.22.12
RUN wget https://go.dev/dl/go1.22.12.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.12.linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/go && \
    rm go1.22.12.linux-amd64.tar.gz

# Set working directory for the build
WORKDIR /go/src/github.com/stratosnet/stratos-chain

COPY go.mod go.sum ./
RUN go mod download

# Add source files
COPY . .
RUN make install


# Final image
FROM ubuntu:20.04

ENV WORK_DIR=/stchaind
ENV RUN_AS_USER=stratos

# Install dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates libgmp-dev curl wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ARG chain_id
ARG moniker
ARG uid=2048
ARG gid=2048

RUN addgroup --gid $gid "$RUN_AS_USER" && \
    useradd --uid $uid --gid $gid --home-dir "$WORK_DIR" --create-home --shell /bin/bash "$RUN_AS_USER"

ENV CHAIN_ID=${chain_id:-DEFAULT}
ENV MONIKER=${moniker:-stratos-test-node}
WORKDIR $WORK_DIR

# Copy over binaries from the build-env
COPY --from=build-env /root/go/bin/stchaind /usr/bin/stchaind
COPY --from=build-env /usr/local/lib/libpbc.so.1.0.0 /usr/local/lib/libpbc.so.1.0.0

RUN cd /usr/local/lib && { ln -s -f libpbc.so.1.0.0 libpbc.so.1 || { rm -f libpbc.so.1 && ln -s libpbc.so.1.0.0 libpbc.so.1; }; } \
  && cd /usr/local/lib && { ln -s -f libpbc.so.1.0.0 libpbc.so || { rm -f libpbc.so && ln -s libpbc.so.1.0.0 libpbc.so; }; } \
  && ldconfig

COPY entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod +x /usr/bin/entrypoint.sh
ENTRYPOINT ["/usr/bin/entrypoint.sh"]
CMD ["stchaind start"]


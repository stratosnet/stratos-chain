FROM golang:1.18-alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3 \
    gmp-dev flex bison

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES
# Install pdc
RUN wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz \
    && tar -xf pbc-0.5.14.tar.gz \
    && cd pbc-0.5.14/ \
    && ./configure \
    && make \
    && make install \
    && ldconfig / \
    && cd .. && rm -rf pbc-0.5.14/ pbc-0.5.14.tar.gz

# Set working directory for the build
WORKDIR /go/src/github.com/stratosnet/stratos-chain

# Add source files
COPY . .
RUN make install


# Final image
FROM alpine:edge

ENV WORK_DIR /stchaind
ENV RUN_AS_USER stratos

# Install ca-certificates
RUN apk add --update ca-certificates gmp-dev

ARG chain_id
ARG moniker
ARG uid=2048
ARG gid=2048

RUN addgroup --gid $gid "$RUN_AS_USER" && \
    adduser -S -G "$RUN_AS_USER" --uid $uid "$RUN_AS_USER" -h "$WORK_DIR"

ENV CHAIN_ID=${chain_id:-DEFAULT}
ENV MONIKER=${moniker:-stratos-node}
WORKDIR $WORK_DIR

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/stchaind /usr/bin/stchaind
COPY --from=build-env /usr/local/lib/libpbc.so.1.0.0 /usr/local/lib/libpbc.so.1.0.0

RUN cd /usr/local/lib && { ln -s -f libpbc.so.1.0.0 libpbc.so.1 || { rm -f libpbc.so.1 && ln -s libpbc.so.1.0.0 libpbc.so.1; }; } \
  && cd /usr/local/lib && { ln -s -f libpbc.so.1.0.0 libpbc.so || { rm -f libpbc.so && ln -s libpbc.so.1.0.0 libpbc.so; }; }

COPY entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod +x /usr/bin/entrypoint.sh
ENTRYPOINT ["/usr/bin/entrypoint.sh"]
CMD ["stchaind start"]

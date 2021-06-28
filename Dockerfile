# Simple usage with a mounted data directory:
# > docker build -t stratos-chain .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.stratos-chain:/stratos-chain/.stratos-chain stratos-chain stratos-chain init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.stratos-chain:/stratos-chain/.stratos-chain stratos-chain stratos-chain start
FROM golang:1.15-alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/stratosnet/stratos-chain

# Add source files
COPY . .

RUN go version

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make install

# Final image
FROM alpine:edge

ENV STRATOS /stchaind

# Install ca-certificates
RUN apk add --update ca-certificates

RUN addgroup stratos && \
    adduser -S -G stratos stratos -h "$STRATOS"

USER stratos

WORKDIR $STRATOS

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/stchaind /usr/bin/stchaind

# Run stchaind by default, omit entrypoint to ease using container with stchaincli
CMD ["stchaind"]

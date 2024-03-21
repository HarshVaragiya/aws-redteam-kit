# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# See docker/README.md for more information.

FROM ubuntu as builder

ENV HOME=/root
ENV KEY_LABEL=thekey
ENV SO_PIN=1234
ENV PIN=1234
ENV LABEL="xks-proxy"

RUN mkdir -p $HOME/aws-kms-xks-proxy
WORKDIR /app/
RUN apt update -y && apt install git softhsm opensc curl build-essential -y && \
    git clone https://github.com/aws-samples/aws-kms-xks-proxy && \
    cp -r /app/aws-kms-xks-proxy/xks-axum/ $HOME/aws-kms-xks-proxy/xks-axum

RUN softhsm2-util --init-token --slot 0 --label $LABEL --so-pin $SO_PIN --pin $PIN
RUN pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so \
                --token-label xks-proxy --login --login-type user \
                --keygen --id F0 --label $KEY_LABEL --key-type aes:32 \
                --pin $PIN

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="$HOME/.cargo/bin:$PATH"

RUN mkdir -p /var/local/xks-proxy/.secret

ENV PROJECT_DIR=$HOME/aws-kms-xks-proxy/xks-axum
RUN cargo build --release --manifest-path=$PROJECT_DIR/Cargo.toml && \
        cp $PROJECT_DIR/target/release/xks-proxy /usr/sbin/xks-proxy

FROM ubuntu
COPY --from=builder /etc/softhsm/ /etc/softhsm/
COPY --from=builder /var/lib/softhsm/ /var/lib/softhsm/
COPY --from=builder /usr/lib/ /usr/lib/
COPY --from=builder /usr/bin/ /usr/bin/
COPY --from=builder /var/local/ /var/local/
COPY --from=builder /usr/sbin/xks-proxy /usr/sbin/xks-proxy
EXPOSE 80
ENV XKS_PROXY_SETTINGS_TOML=/var/local/xks-proxy/.secret/settings.toml \
    RUST_BACKTRACE=1
ENTRYPOINT ["/usr/sbin/xks-proxy"]

# Build Shadow in a stock rust builder container
FROM rust:1.45.2-alpine as shadow
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=America/Los_Angeles
COPY . shadow

# Required dynamic libraries
#
# libdarwinia_shadow.so => /usr/local/lib/libdarwinia_shadow.so (0x7fd26af02000)
# libssl.so.1.1 => /lib/libssl.so.1.1 (0x7fd26ae81000)
# libcrypto.so.1.1 => /lib/libcrypto.so.1.1 (0x7fd26ac02000)
# libsqlite3.so.0 => /usr/lib/libsqlite3.so.0 (0x7fd26ab1a000)
# libc.musl-x86_64.so.1 => /lib/ld64.so.1 (0x7fd26bebb000)
RUN apk add --no-cache gcc go openssl-dev llvm\
    && cd shadow \
    && cargo build --release -vv \
    && mkdir /target \
    && cp target/release/shadow /target/ \
    && cp /usr/local/lib/libdarwinia_shadow.so /target/libdarwinia_shadow.so

# Pull Shadow into a second stage deploy alpine container
FROM alpine:latest
COPY --from=shadow /target /target
RUN mv /target/shadow /usr/local/bin/shadow \
    && mv /target/libdarwinia_shadow.so /usr/local/lib/libdarwinia_shadow.so \
    && cp /lib/libc.musl-x86_64.so.1 /lib/ld64.so.1 \
    && rm -rf /target
EXPOSE 3000
ENTRYPOINT ["shadow"]

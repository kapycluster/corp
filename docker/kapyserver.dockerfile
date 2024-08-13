FROM alpine:latest
RUN apk add --no-cache cni-plugins
ENV PATH=/usr/libexec/cni:$PATH
COPY ./bin/kapyserver /
WORKDIR /
USER 65532:65532

FROM alpine:latest
RUN apk add --no-cache cni-plugins
ENV PATH=/usr/share/cni-plugins/bin:$PATH
COPY ./bin/kapyserver /
WORKDIR /
USER 65532:65532

ENTRYPOINT ["/kapyserver"]

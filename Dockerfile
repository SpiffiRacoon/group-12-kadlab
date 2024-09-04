FROM alpine:latest

RUN apk add --update \
    bash \
    curl \
    iputils-ping \
    netcat-openbsd \
    && rm -rf /var/cache/apk/*

CMD nc -l -p 80

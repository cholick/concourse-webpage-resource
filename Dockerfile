FROM alpine:latest

RUN apk add --no-cache --virtual=build-dependencies wget ca-certificates

ADD assets/ /opt/resource/

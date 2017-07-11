FROM alpine:latest

RUN apk add --no-cache --virtual=build-dependencies wget ca-certificates bash bash-doc bash-completion

ADD assets/ /opt/resource/

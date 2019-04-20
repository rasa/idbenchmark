FROM golang:alpine as builder
LABEL maintainer="Ross Smith II <ross@smithii.com>"

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

COPY . /go/src/github.com/rasa/idbenchmark

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/rasa/idbenchmark \
	&& make test \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM golang:alpine as builder
MAINTAINER Ross Smith II <ross@smithii.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/rasa/go-template

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/rasa/go-template \
	&& make static \
	&& mv go-template /usr/bin/go-template \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/go-template /usr/bin/go-template
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "go-template" ]
CMD [ "--help" ]

FROM golang:alpine as builder
MAINTAINER Ross Smith II <ross@smithii.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/rasa/idbenchmark

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/rasa/idbenchmark \
	&& make static \
	&& mv idbenchmark /usr/bin/idbenchmark \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/idbenchmark /usr/bin/idbenchmark
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "idbenchmark" ]
CMD [ "--help" ]

# build stage
FROM golang:1.13-alpine AS build-env

ADD . /src

RUN apk --update upgrade && \
    apk add gcc libc-dev postgresql-libs sqlite-libs && \
    rm -rf /var/cache/apk/* && \
    # See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN cd /src && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app


# final stage
FROM alpine:3.11.3

ARG WORKDIR="/opt/taiko-web"

WORKDIR ${WORKDIR}

ADD ./static ${WORKDIR}/static
ADD ./templates ${WORKDIR}/templates
COPY ./config/*.yaml ${WORKDIR}/config/
COPY --from=build-env /src/app ${WORKDIR}/

ENTRYPOINT ["./app"]

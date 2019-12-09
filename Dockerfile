ARG GOLANG_VERSION=1.13-alpine3.10
FROM golang:${GOLANG_VERSION} as builder
LABEL maintainer "pgillich ta gmail.com"

ARG BUILD_TAG
ARG BUILD_COMMIT
ARG BUILD_BRANCH
ARG BUILD_TIME

RUN apk --update upgrade && \
    apk add sqlite && \
    apk add gcc && \    
    apk add libc-dev && \
    rm -rf /var/cache/apk/*
# See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY . /src
WORKDIR /src
RUN go build -ldflags "-extldflags '-static' -X $PACKAGE_PATH/configs.BuildTag=$BUILD_TAG -X $PACKAGE_PATH/configs.BuildCommit=$BUILD_COMMIT -X $PACKAGE_PATH/configs.BuildBranch=$BUILD_BRANCH -X $PACKAGE_PATH/configs.BuildTime=$BUILD_TIME"

# Making minimal image

FROM alpine:3.10

ARG RECEIVE_PORT="8088"

RUN apk --update upgrade && \
    apk add sqlite && \
    rm -rf /var/cache/apk/*
# See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=builder "/src/sample-blog" "/sample-blog"

ENTRYPOINT ["/sample-blog", "frontend"]

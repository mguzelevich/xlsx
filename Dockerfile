FROM ubuntu AS build

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -qq
RUN apt-get upgrade -y

RUN apt-get install -yqq --no-install-recommends \
	build-essential \
    git \
    automake

RUN apt-get install -yqq --no-install-recommends \
	ca-certificates \
&& update-ca-certificates 2>/dev/null || true

RUN apt-get install -yqq --no-install-recommends \
	golang sqlite

# # Download Go 1.2.2 and install it to /usr/local/go
# RUN curl -s https://storage.googleapis.com/golang/go1.2.2.linux-amd64.tar.gz| tar -v -C /usr/local -xz
# 
# # Let's people find our Go binaries
# ENV PATH $PATH:/usr/local/go/bin

ENV CGO_ENABLED=1
WORKDIR /src

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags -static" -o /out/xlsxcli .

FROM scratch AS bin

COPY --from=build /out/xlsxcli /

# https://www.docker.com/blog/containerize-your-go-developer-environment-part-1/











#

# FROM golang:1.15.0-alpine AS build
# # FROM alpine
# 
# RUN apk --update upgrade
# RUN apk add gcc
# RUN apk add libc-dev
# 
# RUN apk add sqlite
# # See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
# # RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# 
# # removing apk cache
# RUN rm -rf /var/cache/apk/*
# 
# env CGO_ENABLED=1
# 
# WORKDIR /src
# 
# COPY . .
# 
# RUN go build -o /out/xlsxcli .
# 
# FROM scratch AS bin
# 
# COPY --from=build /out/xlsxcli /
# 
# # https://www.docker.com/blog/containerize-your-go-developer-environment-part-1/

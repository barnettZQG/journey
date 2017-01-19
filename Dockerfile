FROM    alpine:3.4

RUN     apk add --no-cache -U go bash curl gcc g++ libc-dev make musl-dev readline sqlite git 

RUN     mkdir -p /go/src /go/bin /go/pkg
ENV     GOPATH=/go
COPY . /go/src/github.com/barnettzqg/journey
RUN     cd /go/src/github.com/barnettzqg/journey && \
        go get && go build && mkdir /blog &&  cp journey /blog/journey &&  \
        mv content  config.json built-in /blog && rm -rf /go && apk del go git
WORKDIR /blog
COPY docker-entrypoint.sh /bin/docker-entrypoint.sh
RUN chmod 755 /bin/docker-entrypoint.sh
ENTRYPOINT ["/bin/docker-entrypoint.sh"]




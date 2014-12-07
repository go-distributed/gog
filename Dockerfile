FROM golang
ENV GOPATH /root/gopher
ENV GOGPATH github.com/go-distributed/gog
RUN go get $GOGPATH
RUN go install $GOGPATH
ENTRYPOINT $GOPATH/bin/gog -user-message-handler=$GOPATH/src/$GOGPATH/test/query_handler.sh
EXPOSE 8424 9424
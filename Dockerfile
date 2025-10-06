FROM golang:1.25.1

WORKDIR /go/src/github.com/sharkrf/spk-srv
COPY . .

RUN go get -u github.com/jteeuwen/go-bindata/... && go generate && go install

EXPOSE 65200/udp

ENTRYPOINT ["spk-srv"]

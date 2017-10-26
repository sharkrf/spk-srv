FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /go/src/github.com/sharkrf/spk-srv
COPY . .

RUN go-wrapper download github.com/jteeuwen/go-bindata/... && go-wrapper install github.com/jteeuwen/go-bindata/... && go-wrapper download && go generate && go-wrapper install

CMD ["go-wrapper", "run"]

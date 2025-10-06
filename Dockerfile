FROM golang:1.25.1-alpine as builder

RUN apk add --no-cache go-bindata

WORKDIR /app
COPY . .

RUN go generate 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v

FROM alpine

COPY --from=builder /app/spk-srv /app/spk-srv
EXPOSE 65200/udp
ENTRYPOINT ["/app/spk-srv"]

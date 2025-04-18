FROM docker.io/library/golang:tip-alpine3.20 as builder

COPY . /golang

RUN cd /golang/golang && \
    go build .

FROM amd64/alpine:3 as runner

COPY  --from=builder /golang/golang/challenge_fr /go/challenge_fr

MAINTAINER vla2mir@gmail.com

ENTRYPOINT ["/go/challenge_fr"]
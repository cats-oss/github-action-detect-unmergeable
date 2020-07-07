FROM golang:1.14-buster AS builder
ADD . /gh-action-detect-unmergeable/
WORKDIR /gh-action-detect-unmergeable/

RUN ["make", "build", "-j"]

FROM debian:buster-20200607-slim
WORKDIR /root/
COPY --from=builder /gh-action-detect-unmergeable/ghaction_unmergeable_detection .
ENTRYPOINT ["/root/ghaction_unmergeable_detection"]
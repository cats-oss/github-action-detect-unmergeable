FROM golang:1.17.6-buster
ADD . /gh-action-detect-unmergeable/
WORKDIR /gh-action-detect-unmergeable/
RUN ["make", "build", "-j"]

FROM debian:buster-20200607-slim
WORKDIR /root/
COPY --from=builder /gh-action-detect-unmergeable/ghaction_unmergeable_detection .
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
ENTRYPOINT ["/root/ghaction_unmergeable_detection"]

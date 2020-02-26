FROM golang:1.14.0-buster

ADD . /gh-action-detect-unmergeable/

WORKDIR /gh-action-detect-unmergeable/

RUN ["go", "build", "-o", "app"]

ENTRYPOINT ["/gh-action-detect-unmergeable/app"]

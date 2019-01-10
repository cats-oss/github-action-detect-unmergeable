FROM golang:1.11.4-stretch

LABEL "com.github.actions.name"="Detect Unmergeable"
LABEL "com.github.actions.description"="Detect unmergeable pull requests"

ENV GO111MODULE on

ADD . /gh-action-detect-unmergeable/

WORKDIR /gh-action-detect-unmergeable/

RUN ["go", "build", "-o", "app"]

ENTRYPOINT ["/gh-action-detect-unmergeable/app"]

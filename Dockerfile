FROM golang:1.13-buster

LABEL "com.github.actions.name"="Detect Unmergeable"
LABEL "com.github.actions.description"="Detect unmergeable pull requests"

ADD . /gh-action-detect-unmergeable/

WORKDIR /gh-action-detect-unmergeable/

RUN ["go", "build", "-o", "app"]

ENTRYPOINT ["/gh-action-detect-unmergeable/app"]

FROM golang:1.15.0-buster
ADD . /gh-action-detect-unmergeable/
WORKDIR /gh-action-detect-unmergeable/
RUN ["make", "build", "-j"]
ENTRYPOINT ["/gh-action-detect-unmergeable/ghaction_unmergeable_detection"]
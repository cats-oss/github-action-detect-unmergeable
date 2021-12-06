FROM golang:1.17.4-buster
ADD . /gh-action-detect-unmergeable/
WORKDIR /gh-action-detect-unmergeable/
RUN ["make", "build", "-j"]
ENTRYPOINT ["/gh-action-detect-unmergeable/ghaction_unmergeable_detection"]
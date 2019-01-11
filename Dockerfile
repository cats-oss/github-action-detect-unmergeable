FROM node:11.6.0-alpine

LABEL "com.github.actions.name"="Detect Unmergeable"
LABEL "com.github.actions.description"="Detect unmergeable pull requests"

ADD . /gh-action-detect-unmergeable/

WORKDIR /gh-action-detect-unmergeable/
RUN ["yarn", "--production"]

ENTRYPOINT ["node", "/gh-action-detect-unmergeable/src/index.js"]
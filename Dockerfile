FROM node:11.6.0-alpine

LABEL "com.github.actions.name"="Detect Unmergeable"
LABEL "com.github.actions.description"="Detect unmergeable pull requests"

ADD src/ /app/src/
ADD package.json /app/package.json
ADD yarn.lock /app/yarn.lock

WORKDIR /app
RUN ["yarn", "--production"]

ENTRYPOINT ["node", "/app/src/index.js"]

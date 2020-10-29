FROM node:12-slim

WORKDIR /usr/src/app

OPY package*.json ./

RUN yarn install --production

COPY . ./

CMD [ "node", "index.js" ]
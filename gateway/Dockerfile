FROM node:12-slim

WORKDIR /usr/src/app

COPY package.json yarn.lock ./

RUN yarn install --production

COPY . ./

CMD [ "node", "index.js" ]
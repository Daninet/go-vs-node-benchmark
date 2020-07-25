FROM node:14-slim

COPY ./package.json ./package.json
COPY ./package-lock.json ./package-lock.json

RUN npm ci

COPY ./src ./src

ENV NODE_ENV production

CMD ["node", "src/app.js"]

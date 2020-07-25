FROM node:14-slim

COPY ./ ./

RUN npm ci

ENV NODE_ENV production
ENV DATABASE postgres

CMD ["node", "app.js"]

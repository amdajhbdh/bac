FROM node:22-alpine

WORKDIR /app

RUN apk add --no-cache git curl

RUN npm install -g openclaw

ENV NODE_ENV=production

COPY openclaw.json /workspace/openclaw.json

EXPOSE 18789

CMD ["openclaw", "start"]

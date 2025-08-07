FROM node:18-alpine AS base
WORKDIR /app
RUN apk add --no-cache libc6-compat dumb-init

COPY package*.json ./

FROM base AS deps
RUN npm ci --legacy-peer-deps 

FROM base AS dev

COPY --from=deps /app/node_modules ./node_modules

COPY . .

# RUN npm install -g @nestjs/cli nodemon

ENTRYPOINT ["dumb-init", "--"]

EXPOSE 3000

CMD [ "npm", "run", "start:dev" ]
FROM base AS prod
RUN npm ci --only=production

COPY . .

RUN npm run build

RUN addgroup -g 1001 -S nodejs
RUN adduser -S nestjs -u 1001

RUN chown -R nestjs:nodejs /app 

USER nestjs

EXPOSE 3000

CMD [ "node", "dist/main.js" ]

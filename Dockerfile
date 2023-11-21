FROM node:alpine AS node-builder

WORKDIR /backend

COPY tsmodules/package*.json .
RUN npm install

COPY tsmodules/. .
RUN npm run build

FROM registry.heroiclabs.com/heroiclabs/nakama:3.19.0

COPY --from=node-builder /backend/build/*.js /nakama/data/modules/build/
COPY local.yml /nakama/data/

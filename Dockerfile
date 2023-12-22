FROM node:alpine AS node-builder

WORKDIR /backend

COPY tsmodules/package*.json .
RUN npm install

COPY tsmodules/. .
RUN npm run build

FROM heroiclabs/nakama-pluginbuilder:3.20.0 AS go-builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /backend


COPY gomodules/ .

RUN go mod vendor
RUN go build --trimpath --mod=vendor --buildmode=plugin -o ./servicesrpc.so

FROM registry.heroiclabs.com/heroiclabs/nakama:3.20.0

COPY --from=go-builder /backend/*.so /nakama/data/modules/
COPY --from=node-builder /backend/build/*.js /nakama/data/modules/

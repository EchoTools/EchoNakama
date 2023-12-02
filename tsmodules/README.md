# EchoNakama Server Modules

These are the TypeScript modules that are "compiled" using rollup for the gojvm that runs on Nakama. This is a single module, that is combined into the index.js file.

There are specific requirements for the parser to work. Please see https://heroiclabs.com/docs/nakama/server-framework/typescript-runtime/ for details.

## Prequisites

- Docker
- node
- ???

## Setup

`npm install`

## Build

`npm run build`

## Deployment

Deployment is done as part of the Dockerfile that builds nakama. It builds these modules then includes them in the nakama image.
`docker build ..`
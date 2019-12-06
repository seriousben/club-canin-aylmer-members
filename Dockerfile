FROM golang:latest AS getmembers

ENV GOOS=linux GO111MODULE=on

WORKDIR /service

COPY membersfetcher/* ./

RUN go run main.go > members.json

FROM node:13.2.0-alpine as build
WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY package.json /app/package.json
RUN npm install --silent

COPY . /app
COPY --from=getmembers /service/members.json /app/public/.
RUN npm run build

FROM nginx:1.17-alpine
COPY --from=build /app/build /usr/share/nginx/html

FROM golang:1.16.2-alpine AS base
RUN apk add build-base
WORKDIR /app-build
COPY . .

#Build MS
FROM base as build
ARG APP
RUN go get -t ${APP}
RUN go build -o /dist/app ${APP}

#Http MS
FROM alpine AS release
WORKDIR /dist
EXPOSE 8080
COPY --from=build /dist/ /dist/
ENTRYPOINT ./app

#CLI
FROM alpine AS cli
WORKDIR /dist
COPY --from=build /dist/ /dist/
ENTRYPOINT ./app
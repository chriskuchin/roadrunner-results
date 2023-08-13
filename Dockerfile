FROM node:18-alpine AS webpack

ADD . /public
WORKDIR /public
RUN npm install && npm run build

FROM golang:1.20-alpine3.18 AS builder

COPY ${PWD} /app
WORKDIR /app

RUN apk update && apk add --no-cache musl-dev gcc build-base
RUN CGO_ENABLED=1 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin *.go

RUN apk add curl && \
    curl -fsSL -o /go/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 && \
    chmod +x /go/bin/dbmate

FROM alpine:3.18
LABEL MAINTAINER Chris Kuchin <github@kchn.dev>

ENV FRONTEND_FOLDER /home/appuser/app/public/dist
ENV DATABASE_URL sqlite:/rslts/results.db
ENV DB_PATH /rslts/results.db

RUN apk --update add ca-certificates sqlite && \
    mkdir /rslts && rm -rf /var/cache/apk/* && \
    adduser -D appuser && chown appuser /rslts -R


COPY --from=webpack /public/dist /home/appuser/app/public/dist
COPY --from=builder /app/appbin /home/appuser/app/appbin
COPY --from=builder /app/db /home/appuser/app/db
COPY --from=builder /go/bin/dbmate /usr/local/bin/dbmate
COPY ./entrypoint.sh /usr/local/bin/entrypoint

RUN chown -R appuser:appuser /home/appuser/

USER appuser

WORKDIR /home/appuser/app

CMD ["entrypoint"]
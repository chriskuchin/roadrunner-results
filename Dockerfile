FROM node:18-alpine AS webpack

ADD src /public
WORKDIR /public
RUN npm install && npx webpack

FROM golang:1.20-alpine3.17 AS builder

COPY ${PWD} /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin *.go

FROM alpine:3.17
LABEL MAINTAINER Chris Kuchin <github@kchn.dev>

ENV FRONTEND_FOLDER /home/appuser/app/public/dist

RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

RUN adduser -D appuser
USER appuser

COPY --from=webpack /public/dist /home/appuser/app/public/dist
COPY --from=builder /app/appbin /home/appuser/app/appbin

WORKDIR /home/appuser/app

CMD ["./appbin", "server"]
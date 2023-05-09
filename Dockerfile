FROM node:18-alpine AS webpack

ADD . /public
WORKDIR /public
RUN npm install && npm run build

FROM golang:1.20-alpine3.17 AS builder

COPY ${PWD} /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin *.go

RUN go get github.com/amacneil/dbmate && \
    go install github.com/amacneil/dbmate && \
    which dbmate

FROM alpine:3.17
LABEL MAINTAINER Chris Kuchin <github@kchn.dev>

ENV FRONTEND_FOLDER /home/appuser/app/public/dist

RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

RUN adduser -D appuser
USER appuser

COPY --from=webpack /public/dist /home/appuser/app/public/dist
COPY --from=builder /app/appbin /home/appuser/app/appbin
# COPY --from=builder /usr/local/bin /home/appuser/go/

WORKDIR /home/appuser/app

CMD ["./appbin", "server"]
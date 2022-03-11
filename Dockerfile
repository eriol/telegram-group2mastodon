FROM golang:1.17-alpine3.15 AS builder

WORKDIR /app

RUN apk --no-cache add git

COPY . .

RUN go build

FROM alpine:3.15

COPY --from=builder /app/telegram-group2mastodon /bin/telegram-group2mastodon

RUN apk --no-cache add bash

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

ENTRYPOINT /bin/telegram-group2mastodon run

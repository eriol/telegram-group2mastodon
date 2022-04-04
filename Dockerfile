FROM golang:1.17-alpine3.15 AS builder
WORKDIR /app
RUN apk -U --no-cache add git
COPY . .
RUN go build

FROM alpine:3.15
LABEL LastUpdate="2022-04-04"
COPY --from=builder /app/telegram-group2mastodon /bin/telegram-group2mastodon
RUN apk -U --no-cache upgrade
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
ENTRYPOINT /bin/telegram-group2mastodon run

FROM golang:1.24-alpine AS build

COPY . /build
RUN apk update && apk add --no-cache make git
WORKDIR /build

RUN make build

FROM alpine:latest

COPY --from=build /build/bin/telegram-perplexity-bot /usr/bin/telegram-perplexity-bot

RUN apk add --no-cache ca-certificates bash tmux

ENTRYPOINT ["/usr/bin/telegram-perplexity-bot"]

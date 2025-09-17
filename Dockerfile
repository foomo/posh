FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -u 1001 -g 1001 posh

COPY posh /usr/bin/

USER posh
WORKDIR /home/posh

ENTRYPOINT ["posh"]

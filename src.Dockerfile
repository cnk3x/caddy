FROM golang as go

WORKDIR /sources
COPY ./ ./
# RUN ls -lah && exit 1
RUN CGO_ENABLED=0 go build -v -mod=vendor -ldflags '-extldflags "-static"' -o caddy .

FROM alpine as deps

RUN apk --no-cache add upx ca-certificates tzdata
COPY --from=go /sources/caddy /sources/caddy
RUN upx /sources/caddy

FROM scratch

COPY --from=deps /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=deps /sources/caddy /caddy

ENV HOME=/data
WORKDIR /data
VOLUME [ "/data" ]
EXPOSE 80 443

ENTRYPOINT [ "/caddy" ]

CMD [ "run", "-environ", "-watch" ]

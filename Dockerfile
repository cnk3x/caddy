FROM cnk3x/golang:1.18 as builder


ENV GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct

WORKDIR /go/caddy
COPY . .
RUN CGO_ENABLED=0 go build -mod vendor -trimpath -ldflags '-s -w -extldflags -static' -v -o /caddy/caddy ./ && \
    upx /caddy/caddy && \
    echo ":80\nrespond \"hello caddy server!\"" > /caddy/Caddyfile

FROM busybox

COPY --from=builder /caddy/ /

ENV HOME=/data
WORKDIR /data
VOLUME [ "/data" ]
EXPOSE 80 443

ENTRYPOINT [ "/caddy" ]

CMD [ "run", "-environ", "-watch" ]

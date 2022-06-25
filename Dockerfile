FROM cnk3x/golang as builder

RUN echo ":80\nrespond \"hello caddy server!\"" > /rootfs/Caddyfile

ENV GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct

WORKDIR /go/caddy
COPY local_build.sh ./
RUN chmod +x local_build.sh && \
    ./local_build.sh && \
    mv /go/caddy/caddy /rootfs/caddy && \
    upx /caddy/caddy && \
    echo ":80\nrespond \"hello caddy server!\"" > /caddy/Caddyfile

FROM scratch

COPY --from=builder /rootfs/ /

ENV HOME=/data
WORKDIR /data
VOLUME [ "/data" ]
EXPOSE 80 443

ENTRYPOINT [ "/caddy" ]

CMD [ "run", "-environ", "-watch" ]

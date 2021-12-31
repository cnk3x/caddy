# caddy builder

编译 caddy，集成目前官网发布的所有插件

```sh
git clone https://github.com/cnk3x/caddy.git
cd caddy

# 编译到本地
sh ./build.sh

# 编译成镜像
# sh ./build_docker.sh 镜像名称
sh ./build_docker.sh ghcr.io/cnk3x/caddy

```

# caddy builder

编译 caddy，集成目前官网发布的所有插件

```sh
git clone https://github.com/cnk3x/caddy.git
cd caddy

# 本地编译
sh ./build.sh
# 通过 xcaddy 的镜像编译成镜像
# sh ./build_docker.sh 镜像名称
sh ./build_docker.sh ghcr.io/cnk3x/caddy

# 通过本地代码译成镜像
# sh ./build_docker_src.sh 镜像名称
sh ./build_docker.sh ghcr.io/cnk3x/caddy
```

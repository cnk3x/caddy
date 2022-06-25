tag:="cnk3x/caddy"
version:=2.5.1

build:
	@docker build --tag $(tag):latest --tag $(tag):$(version) .
	@docker run --rm $(tag) list-modules
	@docker run --rm $(tag) version
	@echo "you can use \`docker push $(tag)\` to publish this repo"

push:
	@docker push $(tag):latest
	@docker push $(tag):$(version)
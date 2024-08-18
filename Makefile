.PHONY: install

build:
	docker build . -t asteurer/test

push:
	docker push ghcr.io/kube-hack/command-injection

install:
	helm upgrade --install command-injection ./chart

uninstall:
	helm uninstall command-injection
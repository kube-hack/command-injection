.PHONY: install

build:
	docker build . -t ghcr.io/kube-hack/command-injection

push:
	docker push ghcr.io/kube-hack/command-injection

install:
	helm upgrade --install command-injection ./chart

uninstall:
	helm uninstall command-injection
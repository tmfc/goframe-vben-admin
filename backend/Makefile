ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "template-single"
DOCKER_NAME = "template-single"

include ./hack/hack-cli.mk
include ./hack/hack.mk

.PHONY: cover
cover:
	GF_GCFG_FILE=config.test.toml go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html
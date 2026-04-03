GOPATH=$(shell go env GOPATH)
DEPLOY_DIRECTORY=deploy
APPLICATION_NAME=qmanager

.PHONY: initialize-environment

initialize-environment:
	go mod tidy

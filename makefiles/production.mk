.PHONY: build-production-linux build-production-windows build-production-cli

build-production-linux:
	mkdir -p $(DEPLOY_DIRECTORY)/linux
	CGO_ENABLED=1 go build \
		-v \
		-o $(DEPLOY_DIRECTORY)/linux/$(APPLICATION_NAME) \
		main.go
	cp -r config/catalog $(DEPLOY_DIRECTORY)/linux/

build-production-windows:
	mkdir -p $(DEPLOY_DIRECTORY)/windows
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
		-v \
		-o $(DEPLOY_DIRECTORY)/windows/$(APPLICATION_NAME).exe \
		main.go
	cp -r config/catalog $(DEPLOY_DIRECTORY)/windows/

build-production-cli:
	mkdir -p $(DEPLOY_DIRECTORY)/cli
	go build \
		-v \
		-o $(DEPLOY_DIRECTORY)/cli/qmanager-cli \
		cmd/qmanager-cli/main.go

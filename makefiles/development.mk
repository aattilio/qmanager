.PHONY: build-debug purge-artifacts

build-debug:
	$(QT_TOOLS_BIN)/qtdeploy \
		test \
		desktop \
		.

purge-artifacts:
	rm -rf $(DEPLOY_DIRECTORY)
	rm -rf ./src/frontend/moc_*
	rm -rf ./src/frontend/ui_*
	rm -rf ./src/frontend/rcc_*
	rm -f qmanager-cli

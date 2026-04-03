.PHONY: run-all-unit-tests run-backend-tests run-media-integrity-tests

run-all-unit-tests:
	go test \
		-v \
		./src/backend/... \
		./src/core/...

run-backend-tests:
	go test \
		-v \
		./src/backend/hypervisor/... \
		./src/backend/filesystem/...

run-media-integrity-tests:
	go test \
		-v \
		./src/backend/discovery/tests/media_integrity_test.go

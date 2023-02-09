TEST_ENV_FILE = tmp/env

ifneq (,$(wildcard $(TEST_ENV_FILE)))
    include $(TEST_ENV_FILE)
    export
endif

.PHONY: run
.SILENT: run
run:
	go run ./...

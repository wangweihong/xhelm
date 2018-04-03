EXEC_FILE = xhelm
POLICY_EXEC_FILE = ./bin/$(EXEC_FILE)
VERSION_FILE = src/xhelm/version.go
VERSION_MAJOR = 5.4
GO_VERSION = `go version | cut -c 12-`
SRC_VERSION = `git log | head -n 1 | cut -c 8-`

.PHONY: all clean

all: $(POLICY_EXEC_FILE)

$(POLICY_EXEC_FILE):
	go install -v $(EXEC_FILE)

clean:
	@echo "clean ..."
	rm -f $(POLICY_EXEC_FILE)

rebuild : clean all

install:
	@echo "install ..."
	cp $(POLICY_EXEC_FILE) /sbin
	#godep go install costor

test:
	go test -tags "libdm_no_deferred_remove" ./...

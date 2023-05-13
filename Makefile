BUILD_DIR:=.

GIT_IMPORT="github.com/morty-faas/cli/build"
GIT_COMMIT=$$(git rev-parse --short HEAD)
GIT_TAG=$$(git describe --abbrev=0 --tags)

LDFLAGS="-s -w -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT) -X $(GIT_IMPORT).Version=$(GIT_TAG)"

default: build

.PHONY: build
build:
	go build -ldflags $(LDFLAGS) -o $(BUILD_DIR)/morty main.go

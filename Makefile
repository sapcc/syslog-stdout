PKG     = github.com/sapcc/syslog-stdout
PREFIX := /usr

# NOTE: This repo uses Go modules, and uses a synthetic GOPATH at
# $(CURDIR)/.gopath that is only used for the build cache. $GOPATH/src/ is
# empty.
GO            := GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS :=
GO_LDFLAGS    := -s -w

all: build/syslog-stdout

build/syslog-stdout: FORCE
	$(GO) install $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' '$(PKG)'
build/syslog-generator: util/generator.c
	$(CC) -o $@ $<

install: FORCE all
	install -D -m 0755 build/syslog-stdout "$(DESTDIR)$(PREFIX)/bin/syslog-stdout"

clean: FORCE
	rm -rf -- build

vendor: FORCE
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify

.PHONY: FORCE

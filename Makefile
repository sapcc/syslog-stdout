PREFIX = /usr

all: build/syslog-stdout

GO_BUILDFLAGS = -mod vendor
GO_LDFLAGS    =

build/syslog-stdout: FORCE
	go build $(GO_BUILDFLAGS) -ldflags '-s -w $(GO_LDFLAGS)' -o $@ .
build/syslog-generator: util/generator.c
	$(CC) -o $@ $<

install: FORCE all
	install -D -m 0755 build/syslog-stdout "$(DESTDIR)$(PREFIX)/bin/syslog-stdout"

clean: FORCE
	rm -rf -- build

vendor: FORCE
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: FORCE

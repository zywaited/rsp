PHPCFG = /usr/local/php/bin/php-config

PHPEXE := $(shell $(PHPCFG) --php-binary)
PHPDIR := $(shell $(PHPCFG) --prefix)

export PATH := $(PHPDIR)/bin:$(PATH)
export CFLAGS := $(shell $(PHPCFG) --includes) -DLOGLEVEL=error
export LDFLAGS := -L$(shell $(PHPCFG) --prefix)/lib/

export CGO_CFLAGS := $(CFLAGS) $(CGO_CFLAGS)
export CGO_LDFLAGS := $(LDFLAGS) $(CGO_LDFLAGS)


all:
	go build -buildmode=c-shared -o rsp.so *.go

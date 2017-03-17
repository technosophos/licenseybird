.PHONY: build

build:
	go build -o licenseybird *.go

install: build
	install -m 755 ./licenseybird ${DESTDIR}/usr/local/bin/licenseybird

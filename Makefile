GO := go
NAME := purplecat
VERSION := 0.1.0
DIST := $(NAME)-$(VERSION)

all: test build

setup:
	git submodule update --init

update_version:
	@for i in README.md site/content/_index.md; do\
	    sed -e 's!Version-[0-9.]*-yellowgreen!Version-${VERSION}-yellowgreen!g' -e 's!tag/v[0-9.]*!tag/v${VERSION}!g' $$i > a ; mv a $$i; \
	done
	@sed 's/ARG version=".*"/ARG version="${VERSION}"/g' Dockerfile > a ; mv a Dockerfile
	@sed 's/const Version = .*/const Version = "${VERSION}"/g' version.go > a ; mv a version.go
	@echo "Replace version to \"${VERSION}\""	

start:
	make -C site start

stop:
	make -C site start

www:
	make -C site build

test: setup update_version
	$(GO) test -covermode=count -coverprofile=coverage.out $$(go list ./...)

build:
	$(GO) build -o purplecat -v cmd/purplecat/main.go

define _createDist
	mkdir -p dist/$(1)_$(2)/$(DIST)
	GOOS=$1 GOARCH=$2 go build -o dist/$(1)_$(2)/$(DIST)/purplecat$(3) cmd/purplecat/main.go
	cp -r README.md LICENSE dist/$(1)_$(2)/$(DIST)
	tar cfz dist/$(DIST)_$(1)_$(2).tar.gz -C dist/$(1)_$(2) $(DIST)
endef

dist: build
	@$(call _createDist,darwin,amd64,)
	@$(call _createDist,windows,amd64,.exe)
	@$(call _createDist,windows,386,.exe)
	@$(call _createDist,linux,amd64,)
	@$(call _createDist,linux,386,)

clean:
	$(GO) clean
	rm -rf $(NAME)

distclean: clean
	-rm -rf dist

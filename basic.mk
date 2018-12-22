ifneq ($(OS),Windows_NT)
	FIND?=find
	# Set an output prefix, which is the local directory if not specified
	PREFIX?=$(shell pwd)
	STDERR?=/dev/stderr
	SHA256SUM=sha256sum
else
	EXE_EXT?=.exe
	CYGPATH?=$(shell where.exe cygpath.exe | tr "\\\\\\\\" "//")
	FIND?=$(shell where find.exe | grep -E -iv "\\\\System32\\\\" | head -n 1 | tr "\\\\\\\\" "//")
	PREFIX?=$(shell $(CYGPATH) -a -m .)
	MSYS_ROOT?=C:/msys64
	MSYS_ROOT:=$(shell echo '$(MSYS_ROOT)' | tr "\\\\\\\\" "//")
	ifneq ($(wildcard $(MSYS_ROOT)/usr/bin/rm.exe),)
		RM=$(MSYS_ROOT)/usr/bin/rm.exe -f
	endif
	SHA256SUM=$(shell where sha256sum | tr "\\\\\\\\" "//")
endif

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := ${PREFIX}/cross

# Populate version variables
# Add to compile time flags
VERSION := $(shell cat <VERSION.txt)
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif
ifeq ($(GITCOMMIT),)
    GITCOMMIT := ${GITHUB_SHA}
endif
CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

# Set our default go compiler
GO := go

# List the GOOS and GOARCH to build
GOOSARCHES = $(shell grep -v '^\#' .goosarch)

# Set the graph driver as the current graphdriver if not set.
DOCKER_GRAPHDRIVER := $(if $(DOCKER_GRAPHDRIVER),$(DOCKER_GRAPHDRIVER),$(shell docker info 2>&1 | grep "Storage Driver" | sed 's/.*: //'))
export DOCKER_GRAPHDRIVER

# If this session isn't interactive, then we don't want to allocate a
# TTY, which would fail, but if it is interactive, we do want to attach
# so that the user can send e.g. ^C through.
INTERACTIVE := $(shell [ -t 0 ] && echo 1 || echo 0)
ifeq ($(INTERACTIVE), 1)
	DOCKER_FLAGS += -t
endif

.PHONY: build
build: prebuild $(NAME)$(EXE_EXT) ## Builds a dynamic executable or package.

$(NAME)$(EXE_EXT): $(wildcard *.go) $(wildcard */*.go) VERSION.txt
	@echo "+ $@"
	$(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(NAME)$(EXE_EXT) .

.PHONY: static
static: prebuild ## Builds a static executable.
	@echo "+ $@"
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
				-tags "$(BUILDTAGS) static_build" \
				${GO_LDFLAGS_STATIC} -o $(NAME)$(EXE_EXT) .

all: clean build fmt lint test staticcheck vet install ## Runs a clean, build, fmt, lint, test, staticcheck, vet and install.

.PHONY: fmt
fmt: ## Verifies all files have been `gofmt`ed.
	@echo "+ $@"
	@gofmt -s -l . | grep -E -v '(\.pb\.go:|vendor)' | tee $(STDERR)

.PHONY: lint
lint: ## Verifies `golint` passes.
	@echo "+ $@"
	@golint ./... | grep -E -v '(\.pb\.go:|vendor)' | tee $(STDERR)

.PHONY: test
test: prebuild ## Runs the go tests.
	@echo "+ $@"
	@$(GO) test -v -tags "$(BUILDTAGS) cgo" $(shell $(GO) list ./... | grep -v vendor)

.PHONY: vet
vet: ## Verifies `go vet` passes.
	@echo "+ $@"
	@$(GO) vet $(shell $(GO) list ./... | grep -E -v '(\.pb\.go:|vendor)') | tee $(STDERR)

.PHONY: staticcheck
staticcheck: ## Verifies `staticcheck` passes.
	@echo "+ $@"
	@staticcheck $(shell $(GO) list ./... | grep -E -v '(\.pb\.go:|vendor)') | tee $(STDERR)

.PHONY: cover
cover: prebuild ## Runs go test with coverage.
	@echo "" > coverage.txt
	@for d in $(shell $(GO) list ./... | grep -v vendor); do \
		$(GO) test -race -coverprofile=profile.out -covermode=atomic "$$d"; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done;

.PHONY: install
install: prebuild ## Installs the executable or package.
	@echo "+ $@"
	$(GO) install -a -tags "$(BUILDTAGS)" ${GO_LDFLAGS} .

define buildpretty
mkdir -p $(BUILDDIR)/$(1)/$(2);
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
	 -o $(BUILDDIR)/$(1)/$(2)/$(NAME)$(if $(findstring windows,$(1)),.exe) \
	 -a -tags "$(BUILDTAGS) static_build netgo" \
	 -installsuffix netgo ${GO_LDFLAGS_STATIC} .;
md5sum $(BUILDDIR)/$(1)/$(2)/$(NAME)$(if $(findstring windows,$(1)),.exe) > $(BUILDDIR)/$(1)/$(2)/$(NAME)$(if $(findstring windows,$(1)),.exe).md5;
$(SHA256SUM) $(BUILDDIR)/$(1)/$(2)/$(NAME)$(if $(findstring windows,$(1)),.exe) > $(BUILDDIR)/$(1)/$(2)/$(NAME)$(if $(findstring windows,$(1)),.exe).sha256;
endef

.PHONY: cross
cross: *.go VERSION.txt prebuild ## Builds the cross-compiled binaries, creating a clean directory structure (eg. GOOS/GOARCH/binary).
	@echo "+ $@"
	$(foreach GOOSARCH,$(GOOSARCHES), $(call buildpretty,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH))))

define buildrelease
echo -n;
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
	 -o $(BUILDDIR)/$(NAME)-$(1)-$(2)$(if $(findstring windows,$(1)),.exe) \
	 -a -tags "$(BUILDTAGS) static_build netgo" \
	 -installsuffix netgo ${GO_LDFLAGS_STATIC} .;
md5sum $(BUILDDIR)/$(NAME)-$(1)-$(2)$(if $(findstring windows,$(1)),.exe) > $(BUILDDIR)/$(NAME)-$(1)-$(2)$(if $(findstring windows,$(1)),.exe).md5;
$(SHA256SUM) $(BUILDDIR)/$(NAME)-$(1)-$(2)$(if $(findstring windows,$(1)),.exe) > $(BUILDDIR)/$(NAME)-$(1)-$(2)$(if $(findstring windows,$(1)),.exe).sha256;
endef

.PHONY: release
release: *.go VERSION.txt prebuild ## Builds the cross-compiled binaries, naming them in such a way for release (eg. binary-GOOS-GOARCH).
	@echo "+ $@"
	$(foreach GOOSARCH,$(GOOSARCHES), $(call buildrelease,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH))))

.PHONY: bump-version
BUMP := patch
bump-version: ## Bump the version in the version file. Set BUMP to [ patch | major | minor ].
	@$(GO) get -u github.com/jessfraz/junk/sembump # update sembump tool
	$(eval NEW_VERSION = $(shell sembump --kind $(BUMP) $(VERSION)))
	@echo "Bumping VERSION.txt from $(VERSION) to $(NEW_VERSION)"
	echo $(NEW_VERSION) > VERSION.txt
	@echo "Updating links to download binaries in README.md"
	sed -i s/$(VERSION)/$(NEW_VERSION)/g README.md
	git add VERSION.txt README.md
	git commit -vsam "Bump version to $(NEW_VERSION)"
	@echo "Run make tag to create and push the tag for new version $(NEW_VERSION)"

.PHONY: tag
tag: ## Create a new git tag to prepare to build a release.
	git tag -sa $(VERSION) -m "$(VERSION)"
	@echo "Run git push origin $(VERSION) to push your new tag to GitHub and trigger a travis build."

REGISTRY := r.smithii.com
.PHONY: image
image: ## Create the docker image from the Dockerfile.
	@docker build --rm --force-rm -t $(REGISTRY)/$(NAME) .

.PHONY: image-dev
image-dev:
	@docker build --rm --force-rm -f Dockerfile.dev -t $(REGISTRY)/$(NAME):dev .

.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@

.PHONY: vendor
vendor: ## Updates the vendoring directory.
	@$(RM) go.sum
	@test -d vendor && $(RM) -r vendor || true
	GO111MODULE=on $(GO) mod init || true
	GO111MODULE=on $(GO) mod tidy
	GO111MODULE=on $(GO) mod vendor
	@$(RM) Gopkg.toml Gopkg.lock

.PHONY: clean
clean: ## Cleanup any build binaries or packages.
	@echo "+ $@"
	$(RM) $(NAME)$(EXE_EXT)
	$test -d $(BUILDDIR) && (RM) -r $(BUILDDIR)

.PHONY: gofmt
gofmt: ## Format all .go files via `gofmt -s` (simplify)
	@echo "+ $@"
	@gofmt -s -l . | grep -E -v '(\.pb\.go:|vendor)' || true
	$(FIND) . -iname '*.go' ! -ipath './vendor/*' | xargs gofmt -s -w

.PHONY: mailmap
mailmap: ## Generate committer list to add to .mailmap
	git log --format='%aN <%aE>' | sort -uf | grep -v -E '\w{8}-\w{4}-\w{4}-\w{4}-\w{10}' | grep -iv 'nobody@localhost'

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/^[^:]*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

check_defined = \
    $(strip $(foreach 1,$1, \
	$(call __check_defined,$1,$(strip $(value 2)))))

__check_defined = \
    $(if $(value $1),, \
    $(error Undefined $1$(if $2, ($2))$(if $(value @), \
    required by target `$@')))

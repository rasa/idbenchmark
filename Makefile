# Setup name variables for the package/tool
NAME := idbenchmark
PKG := github.com/rasa/$(NAME)

CGO_ENABLED := 0

# Set any default go build tags.
BUILDTAGS :=

include basic.mk

ifneq ("$(wildcard local.mk)", "")
include local.mk
endif

.PHONY: prebuild
prebuild:

TZ := UTC
export TZ

# Reading dotenv
DOTENV ?=.env
ifneq (, $(DOTENV))
ifneq (,$(wildcard ./$(DOTENV)))
	include $(DOTENV)
	export
endif
endif

# GO Options
SRC        	:= $(shell find . -type f -name '*.go' -print)
TAGS      	?=
TEST_TAGS   ?=
GOFLAGS    	?=
EXT_LDFLAGS ?=
LDFLAGS      = -w -s $(EXT_LDFLAGS)


.PHONY: build
build:
	go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./...


.PHONY: run
run:
	go run $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' main.go


.PHONY: all build version

ECHO ?= echo
REVISION_NAME=REVISION

all: build version

build:
	# Build your app here: it can be an archive with your binary and revision file

version:
	${ECHO} $(shell git log -n 1 --pretty=format:"%H") > ${REVISION_NAME}

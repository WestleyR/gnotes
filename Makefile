#
#  Makefile
#
# Created by WestleyR on 2022-02-20
# Source code: https://github.com/WestleyR/gnotes
#
# Copyright (c) 2022 WestleyR. All rights reserved.
# This software is licensed under a BSD 3-Clause Clear License.
# Consult the LICENSE file that came with this software regarding
# your rights to distribute this software.
#

# Only cli backend interface
TARGET_BACKEND = gnotes-backend

# CLI interface
TARGET_CLI = gnotes

GO = go
GOFLAGS = -ldflags -w

SRC = $(shell find . -name '*.go')

all: $(TARGET_CLI)

$(TARGET_CLI): $(SRC)
	$(GO) build $(GOFLAGS) -o $(TARGET_CLI) ./cmd/cli

generate: $(SRC)
	go build -buildmode c-shared -o bridge-c/gnotes-bridge.so bridge-go/*.go

build-c: generate
	gcc example-c/main.c bridge-c/gnotes-bridge.so

IOS_OUT = ios-lib

ios-arm64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=arm64 \
	SDK=iphoneos \
	CC=$(PWD)/clangwrap.sh \
	CGO_CFLAGS="-fembed-bitcode" \
	go build -buildmode=c-archive -tags ios -o $(IOS_OUT)/arm64.a bridge-go/*.go

ios-x86_64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	SDK=iphonesimulator \
	CC=$(PWD)/clangwrap.sh \
	go build -buildmode=c-archive -tags ios -o $(IOS_OUT)/x86_64.a bridge-go/*.go

ios: ios-arm64 ios-x86_64
	lipo $(IOS_OUT)/x86_64.a $(IOS_OUT)/arm64.a -create -output $(IOS_OUT)/gnotes.a
	cp $(IOS_OUT)/arm64.h $(IOS_OUT)/gnotes.h

clean:
	rm -f $(TARGET_GNOTES) $(TARGET_CLI)


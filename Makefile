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
	gcc -g -Wall example-c/main.c bridge-c/gnotes-bridge.so

clean:
	rm -f $(TARGET_GNOTES) $(TARGET_CLI)


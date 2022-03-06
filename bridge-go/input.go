//
// Created by WestleyR <westleyr@nym.hush.com> on 2022-03-06
// Source code: https://github.com/WestleyR/gnotes
//
// Copyright (c) 2022 WestleyR. All rights reserved.
// This software is licensed under a BSD 3-Clause Clear License.
// Consult the LICENSE file that came with this software regarding
// your rights to distribute this software.
//

package main

import (
	"log"
	"strings"
)

func parseInput(input string) map[string]string {
	ret := make(map[string]string)

	if input == "" {
		return ret
	}

	for _, opts := range strings.Split(input, " ") {
		op := strings.Split(opts, "=")

		if len(op) != 2 {
			log.Fatalf("invalid input: %s", op)
		}

		ret[op[0]] = op[1]
	}

	return ret
}

func isBool(s string) bool {
	switch s {
	case "true", "t", "yes", "1", "TRUE":
		return true
	}
	return false
}

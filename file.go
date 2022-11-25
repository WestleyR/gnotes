//
//  files.go - https://github.com/WestleyR/gnotes
//
// Created by WestleyR <westleyr@nym.hush.com> on 2022-02-20
// Source code: https://github.com/WestleyR/gnotes
//
// Copyright (c) 2022 WestleyR. All rights reserved.
// This software is licensed under a BSD 3-Clause Clear License.
// Consult the LICENSE file that came with this software regarding
// your rights to distribute this software.
//

package gnotes

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func Sha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

func Sha1File(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	sha1 := Sha1(string(b))

	return sha1, nil
}

func formatBytes(bytes int64) string {
	const unit = 1023

	div := int64(unit)
	exp := 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "kMGTPE"[exp])
}

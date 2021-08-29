//
//  config.go - https://github.com/WestleyR/gnotes
//  gnotes - CLI based S3 syncing note app
//
// Created by WestleyR <westleyr@nym.hush.com> on 2021-08-28
// Source code: https://github.com/WestleyR/gnotes
//
// Copyright (c) 2021 WestleyR. All rights reserved.
// This software is licensed under a BSD 3-Clause Clear License.
// Consult the LICENSE file that came with this software regarding
// your rights to distribute this software.
//

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vaughan0/go-ini"
)

const appID = "wst.gnotes"

func getFileFromConfig(file string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	configFile := filepath.Join(home, ".config", appID)

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		err = os.MkdirAll(configFile, 0700)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create: %s: %s\n", configFile, err)
			return ""
		}
	}

	configFile = filepath.Join(configFile, file)

	f, err := os.OpenFile(configFile, os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open: %s: %s\n", configFile, err)
		return ""
	}
	f.Close()

	return configFile
}

func getConfigFile() string {
	return getFileFromConfig("config.ini")
}

func getLocalSaveFile() string {
	return getFileFromConfig("notes.json")
}

func getConfigValue(header, key string, panicIfEmpty bool) string {
	file, err := ini.LoadFile(getConfigFile())
	if err != nil {
		panic(err)
	}

	value, ok := file.Get(header, key)

	if panicIfEmpty && !ok {
		log.Fatalf("%s/%s not set in %s\n", header, key, getConfigFile())
	}

	return value
}

func getUseS3() bool {
	return getConfigValue("s3", "active", false) == "true"
}

func getS3FileName() string {
	return getConfigValue("s3", "savefile", true)
}

func getS3AccessKey() string {
	return getConfigValue("s3", "accesskey", true)
}

func getS3SecretKey() string {
	return getConfigValue("s3", "secretkey", true)
}

func getS3Region() string {
	return getConfigValue("s3", "region", true)
}

func getS3Bucket() string {
	return getConfigValue("s3", "bucket", true)
}

func getS3Endpoint() string {
	return getConfigValue("s3", "endpoint", true)
}

func getEditor() string {
	return getConfigValue("settings", "editor", true)
}
